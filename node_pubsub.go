// package dsn 定义了分布式存储网络的核心功能
package dsn

// 导入所需的包
import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

// 定义常量
const (
	DefaultLibp2pPubSubMaxMessageSize = 50 << 20              // 默认的libp2p pubsub最大消息大小，设置为50MB
	DefaultPubsubProtocol             = "/dep2p/pubsub/1.0.0" // 默认的pubsub协议版本
)

// PubSubMsgHandler 定义了处理其他节点发布消息的函数类型
type PubSubMsgHandler func(*Message)

// DSN 表示分布式存储网络的主要结构
type DSN struct {
	ctx              context.Context          // 上下文，用于控制goroutine的生命周期
	cancel           context.CancelFunc       // 取消函数，用于取消上下文
	host             host.Host                // libp2p主机，代表网络中的一个节点
	pubsub           *PubSub                  // PubSub实例，用于发布订阅功能
	topicLock        sync.Mutex               // 主题锁，用于保护主题映射的并发访问
	topicMap         map[string]*Topic        // 主题映射，存储所有已创建的主题
	startUp          int32                    // 启动状态，用原子操作来保证线程安全
	subscribedTopics map[string]*Subscription // 已订阅的主题，存储所有当前节点订阅的主题
	subscribeLock    sync.Mutex               // 订阅锁，用于保护订阅操作的并发访问
	// discoveryService discovery.Discovery      // 发现服务，用于在网络中发现其他节点
}

// NewDSN 创建并返回一个新的 DSN 实例
// 参数:
//   - ctx: 上下文，用于控制DSN实例的生命周期
//   - host: libp2p主机，代表当前节点
//   - opts: 节点选项，用于自定义DSN的行为
//
// 返回:
//   - *DSN: 新创建的DSN实例
//   - error: 如果创建过程中出现错误，返回相应的错误信息
func NewDSN(ctx context.Context, host host.Host, opts ...NodeOption) (*DSN, error) {
	// 初始化选项，应用用户提供的自定义选项
	options := DefaultOptions()
	if err := options.ApplyOptions(opts...); err != nil {
		return nil, err
	}

	// 创建可取消的上下文，用于控制DSN及其子组件的生命周期
	ctx, cancel := context.WithCancel(ctx)

	// 创建发现服务，用于在网络中发现其他节点
	// discoveryService, err := createDiscoveryService(ctx, host)
	// if err != nil {
	// 	cancel()
	// 	return nil, err
	// }

	// 初始化DSN实例，设置各个字段的初始值
	dsn := &DSN{
		ctx:              ctx,
		cancel:           cancel,
		host:             host,
		topicMap:         make(map[string]*Topic),
		startUp:          0,
		subscribedTopics: make(map[string]*Subscription),
		// discoveryService: discoveryService,
	}

	// 启动 PubSub 服务
	if err := dsn.startPubSub(options); err != nil {
		cancel()
		return nil, err
	}

	return dsn, nil
}

// startPubSub 启动 PubSub 服务
// 参数:
//   - options: 包含PubSub配置的选项
//
// 返回:
//   - error: 如果启动过程中出现错误，返回相应的错误信息
func (dsn *DSN) startPubSub(options *Options) error {
	// 创建JSON追踪器，用于记录PubSub的活动
	tracer, err := NewJSONTracer("/tmp/trace.out.json")
	if err != nil {
		logrus.Errorf("[DSN] 创建JSON追踪器失败: %v", err)
		return err
	}

	// 设置GossipSub参数，根据提供的选项自定义行为
	params := DefaultGossipSubParams()
	params.D = options.D
	params.Dlo = options.Dlo
	params.HeartbeatInterval = options.HeartbeatInterval
	params.IWantFollowupTime = options.FollowupTime

	// 配置PubSub选项，包括追踪器、对等交换、GossipSub参数等
	pubsubOpts := []Option{
		WithEventTracer(tracer),
		WithPeerExchange(true),
		WithGossipSubParams(params),
		// WithDiscovery(dsn.discoveryService),
		WithFloodPublish(true),
		WithMessageSigning(options.SignMessages),
		WithStrictSignatureVerification(options.ValidateMessages),
		WithPeerScore(
			&PeerScoreParams{
				TopicScoreCap:    100,
				AppSpecificScore: func(peer.ID) float64 { return 0 },
				DecayInterval:    time.Second,
				DecayToZero:      0.01,
			},
			&PeerScoreThresholds{
				GossipThreshold:             -1,
				PublishThreshold:            -2,
				GraylistThreshold:           -3,
				OpportunisticGraftThreshold: 1,
			},
		),
		WithMaxMessageSize(DefaultLibp2pPubSubMaxMessageSize),
	}

	// 创建GossipSub实例
	ps, err := NewGossipSub(dsn.ctx, dsn.host, pubsubOpts...)
	if err != nil {
		return err
	}

	// 设置PubSub实例和启动状态
	dsn.pubsub = ps
	atomic.StoreInt32(&dsn.startUp, 2)
	logrus.Info("[DSN] gossip-sub 服务已启动")
	return nil
}

// GetTopic 根据给定的名称获取一个 topic
// 参数:
//   - name: 主题名称
//
// 返回:
//   - *Topic: 获取或创建的主题实例
//   - error: 如果获取或创建过程中出现错误，返回相应的错误信息
func (dsn *DSN) GetTopic(name string) (*Topic, error) {
	// 检查PubSub是否已启动
	if atomic.LoadInt32(&dsn.startUp) < 2 {
		return nil, fmt.Errorf("libp2p gossip-sub 未运行")
	}

	// 加锁以保护主题映射的并发访问
	dsn.topicLock.Lock()
	defer dsn.topicLock.Unlock()

	// 检查主题是否已存在，如果不存在则创建
	t, ok := dsn.topicMap[name]
	if !ok || t == nil {
		topic, err := dsn.pubsub.Join(name)
		if err != nil {
			return nil, err
		}
		dsn.topicMap[name] = topic
		t = topic
	}
	return t, nil
}

// Subscribe 订阅一个 topic
// 参数:
//   - topic: 主题名称
//   - subscribe: 是否实际进行订阅操作
//
// 返回:
//   - *Subscription: 如果subscribe为true，返回订阅实例；否则返回nil
//   - error: 如果订阅过程中出现错误，返回相应的错误信息
func (dsn *DSN) Subscribe(topic string, subscribe bool) (*Subscription, error) {
	// 如果主题为空，使用默认主题
	if topic == "" {
		topic = DefaultPubsubProtocol
	}
	// 获取主题
	t, err := dsn.GetTopic(topic)
	if err != nil {
		return nil, err
	}
	logrus.Infof("[DSN] gossip-sub 订阅 topic[%s]", topic)

	// 如果需要订阅，则返回订阅实例
	if subscribe {
		return t.Subscribe()
	}
	return nil, nil
}

// Publish 向 topic 发布一条消息
// 参数:
//   - topic: 主题名称
//   - data: 要发布的消息数据
//
// 返回:
//   - error: 如果发布过程中出现错误，返回相应的错误信息
func (dsn *DSN) Publish(topic string, data []byte) error {
	// 获取主题
	t, err := dsn.GetTopic(topic)
	if err != nil {
		return err
	}
	// 发布消息
	return t.Publish(dsn.ctx, data)
}

// IsSubscribed 检查给定的主题是否已经订阅
// 参数:
//   - topic: 主题名称
//
// 返回:
//   - bool: 如果主题已订阅返回true，否则返回false
func (dsn *DSN) IsSubscribed(topic string) bool {
	_, ok := dsn.subscribedTopics[topic]
	return ok
}

// BroadcastWithTopic 将消息广播到给定主题
// 参数:
//   - topic: 主题名称
//   - data: 要广播的消息数据
//
// 返回:
//   - error: 如果广播过程中出现错误，返回相应的错误信息
func (dsn *DSN) BroadcastWithTopic(topic string, data []byte) error {
	// 检查主题是否已订阅
	_, ok := dsn.subscribedTopics[topic]
	if !ok {
		return fmt.Errorf("主题未订阅")
	}
	// 发布消息
	return dsn.Publish(topic, data)
}

// CancelSubscribeWithTopic 取消订阅给定主题
// 参数:
//   - topic: 要取消订阅的主题名称
//
// 返回:
//   - error: 如果取消订阅过程中出现错误，返回相应的错误信息
func (dsn *DSN) CancelSubscribeWithTopic(topic string) error {
	// 检查PubSub是否已启动
	if atomic.LoadInt32(&dsn.startUp) < 2 {
		return fmt.Errorf("libp2p gossip-sub 未运行")
	}

	// 加锁以保护订阅映射的并发访问
	dsn.subscribeLock.Lock()
	defer dsn.subscribeLock.Unlock()

	// 取消订阅并从映射中删除
	if topicSub, ok := dsn.subscribedTopics[topic]; ok {
		if topicSub != nil {
			topicSub.Cancel()
		}
		delete(dsn.subscribedTopics, topic)
	}
	return nil
}

// CancelPubsubWithTopic 取消给定名字的订阅
// 参数:
//   - name: 要取消的主题名称
//
// 返回:
//   - error: 如果取消过程中出现错误，返回相应的错误信息
func (dsn *DSN) CancelPubsubWithTopic(name string) error {
	// 检查PubSub是否已启动
	if atomic.LoadInt32(&dsn.startUp) < 2 {
		return fmt.Errorf("libp2p gossip-sub 未运行")
	}

	// 加锁以保护主题映射的并发访问
	dsn.topicLock.Lock()
	defer dsn.topicLock.Unlock()

	// 关闭主题并从映射中删除
	if topic, ok := dsn.topicMap[name]; ok {
		if err := topic.Close(); err != nil {
			return err
		}
		delete(dsn.topicMap, name)
	}
	return nil
}

// SubscribeWithTopic 订阅给定主题，并使用给定的订阅消息处理函数
// 参数:
//   - topic: 要订阅的主题名称
//   - handler: 用于处理接收到的消息的函数
//   - subscribe: 是否实际进行订阅操作
//
// 返回:
//   - error: 如果订阅过程中出现错误，返回相应的错误信息
func (dsn *DSN) SubscribeWithTopic(topic string, handler PubSubMsgHandler, subscribe bool) error {
	// 加锁以保护订阅映射的并发访问
	dsn.subscribeLock.Lock()
	defer dsn.subscribeLock.Unlock()

	// 检查主题是否已订阅
	if dsn.IsSubscribed(topic) {
		return fmt.Errorf("主题已订阅")
	}

	// 订阅主题
	topicSub, err := dsn.Subscribe(topic, subscribe)
	if err != nil {
		return err
	}

	// 将订阅添加到映射中
	dsn.subscribedTopics[topic] = topicSub

	// 如果订阅成功，启动消息处理循环
	if topicSub != nil {
		go dsn.topicSubLoop(topicSub, handler)
	}

	return nil
}

// // topicSubLoop 用于处理订阅主题的消息循环
// // 参数:
// //   - topicSub: 订阅实例
// //   - handler: 用于处理接收到的消息的函数
// func (dsn *DSN) topicSubLoop(topicSub *Subscription, handler PubSubMsgHandler) {
// 	for {
// 		// 获取下一条消息
// 		message, err := topicSub.Next(dsn.ctx)

// 		if err != nil {
// 			if err.Error() == "subscription cancelled" {
// 				logrus.Warn("[DSN] 订阅已取消：", err)
// 				break
// 			}
// 			logrus.Errorf("[DSN] 订阅下一个消息失败：%s", err.Error())
// 		}

// 		if message == nil {
// 			return
// 		}

// 		// 忽略自己发送的消息
// 		if message.GetFrom() == dsn.host.ID() {
// 			continue
// 		}

// 		// 解析消息
// 		var request Message
// 		if err := request.Unmarshal(message.Data); err != nil {
// 			return
// 		}

// 		if len(request.From) == 0 {
// 			return
// 		}

// 		pid, err := peer.IDFromBytes(request.From) // 从消息中提取 peer.ID
// 		if err != nil {
// 			return
// 		}

// 		// 检查接收者是否为自己
// 		if pid.String() != dsn.host.ID().String() {
// 			continue
// 		}

// 		// 处理消息
// 		handler(&request)
// 	}
// }

// topicSubLoop 用于处理订阅主题的消息循环
// 参数:
//   - topicSub: 订阅实例
//   - handler: 用于处理接收到的消息的函数
func (dsn *DSN) topicSubLoop(topicSub *Subscription, handler PubSubMsgHandler) {
	for {
		// 获取下一条消息
		message, err := topicSub.Next(dsn.ctx)

		if err != nil {
			if err.Error() == "subscription cancelled" {
				logrus.Warn("[DSN] 订阅已取消：", err)
				break
			}
			logrus.Errorf("[DSN] 订阅下一个消息失败：%s", err.Error())
		}

		if message == nil {
			return
		}

		// 忽略自己发送的消息
		if message.GetFrom() == dsn.host.ID() {
			continue
		}

		// 解析消息
		// var request Message
		// fmt.Println(len(message.Data))
		// if err := request.Unmarshal(message); err != nil {
		// 	return
		// }

		if len(message.From) == 0 {
			return
		}

		pid, err := peer.IDFromBytes(message.From) // 从消息中提取 peer.ID
		if err != nil {

			return
		}

		// 检查接收者是否为自己
		if pid.String() == dsn.host.ID().String() {
			continue
		}

		// 处理消息
		handler(message)
	}
}

// Pubsub 返回 PubSub 实例
// 返回:
//   - *PubSub: 当前DSN实例使用的PubSub实例
func (dsn *DSN) Pubsub() *PubSub {
	return dsn.pubsub
}

// ListPeers 返回我们在给定主题中连接到的对等点列表
// 参数:
//   - topic: 主题名称
//
// 返回:
//   - []peer.ID: 与给定主题相关的对等点ID列表
func (dsn *DSN) ListPeers(topic string) []peer.ID {
	return dsn.pubsub.ListPeers(topic)
}
