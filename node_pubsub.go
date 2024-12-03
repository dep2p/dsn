// package dsn 定义了分布式存储网络的核心功能
package dsn

// 导入所需的包
import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

// 定义常量
const (
	DefaultPubsubProtocol = "/dep2p/pubsub/1.0.0" // 默认的pubsub协议版本
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
// 该函数负责初始化和配置发布订阅系统，支持多种发布订阅模式和可选的配置加载
//
// 参数:
//   - options: 包含PubSub配置的选项，包括:
//   - LoadConfig: 是否加载详细配置
//   - PubSubMode: 发布订阅模式（GossipSub/FloodSub/RandomSub）
//   - MaxMessageSize: 最大消息大小
//   - 其他详细配置项（当LoadConfig为true时使用）
//
// 返回:
//   - error: 如果启动过程中出现错误，返回相应的错误信息
func (dsn *DSN) startPubSub(options *Options) error {
	// 创建基本的 PubSub 选项数组，用于存储所有 PubSub 配置选项
	var pubsubOpts []Option

	// 根据 LoadConfig 选项决定是否加载详细配置
	if options.GetLoadConfig() {
		// 创建JSON追踪器，用于记录和追踪PubSub网络中的事件和消息流
		// 获取系统的临时目录
		tempDir := os.TempDir()
		// 构建临时文件路径
		traceFilePath := filepath.Join(tempDir, "trace.out.json")
		tracer, err := NewJSONTracer(traceFilePath)
		if err != nil {
			logrus.Errorf("[DSN] 创建JSON追踪器失败: %v", err)
			return err
		}

		// 所有模式通用的基础配置选项
		baseOpts := []Option{
			WithEventTracer(tracer),                                   // 启用事件追踪，用于调试和监控
			WithMessageSigning(options.SignMessages),                  // 配置消息签名选项
			WithStrictSignatureVerification(options.ValidateMessages), // 配置签名验证严格程度
			WithMaxMessageSize(options.MaxMessageSize),                // 设置最大消息大小限制
		}

		// 根据不同的发布订阅模式添加特定的配置选项
		switch options.GetPubSubMode() {
		case GossipSub:
			// GossipSub 特定的参数配置
			// 设置 GossipSub 参数
			params := DefaultGossipSubParams()

			// 设置每个节点维护的对等点数量，最小为2
			params.D = int(math.Max(2, float64(options.D)))

			// 设置对等点数量的最小阈值，最小为1
			params.Dlo = int(math.Max(1, float64(options.Dlo)))

			// 设置心跳间隔，不超过1秒
			if options.HeartbeatInterval > time.Second {
				params.HeartbeatInterval = time.Second
			} else {
				params.HeartbeatInterval = options.HeartbeatInterval
			}

			// 设置消息请求的跟进时间，不超过2秒
			if options.FollowupTime > 2*time.Second {
				params.IWantFollowupTime = 2 * time.Second
			} else {
				params.IWantFollowupTime = options.FollowupTime
			}

			// GossipSub 特定的选项
			gossipOpts := []Option{
				WithPeerExchange(true),      // 启用对等节点交换
				WithGossipSubParams(params), // 设置 GossipSub 协议参数
				WithFloodPublish(true),      // 启用洪泛式消息发布
				// 配置节点评分系统，仅用于 GossipSub
				WithPeerScore(
					&PeerScoreParams{
						TopicScoreCap:    100,                                // 主题分数上限
						AppSpecificScore: func(peer.ID) float64 { return 0 }, // 应用特定的评分函数,这里默认返回0
						DecayInterval:    time.Second,                        // 分数衰减间隔
						DecayToZero:      0.01,                               // 衰减到0的速率
					},
					&PeerScoreThresholds{
						GossipThreshold:             -1, // Gossip消息传播的阈值
						PublishThreshold:            -2, // 发布消息的阈值
						GraylistThreshold:           -3, // 灰名单阈值
						OpportunisticGraftThreshold: 1,  // 机会性嫁接的阈值
					},
				),
			}
			pubsubOpts = append(baseOpts, gossipOpts...)

		case FloodSub:
			// FloodSub 使用基础配置即可
			pubsubOpts = baseOpts

		case RandomSub:
			// RandomSub 特定的配置（如果有的话）
			pubsubOpts = baseOpts
			// 可以在这里添加 RandomSub 特定的选项
		}
	} else {
		// 如果不加载详细配置，只使用基本的消息大小限制
		pubsubOpts = []Option{
			WithMaxMessageSize(options.MaxMessageSize),
		}
	}

	// 根据配置的 PubSubMode 创建相应的发布订阅实例
	var ps *PubSub
	var err error

	switch options.GetPubSubMode() {
	case FloodSub:
		ps, err = NewFloodSub(dsn.ctx, dsn.host, pubsubOpts...)
		if err == nil {
			logrus.Info("[DSN] flood-sub 服务已启动")
		}
	case RandomSub:
		return fmt.Errorf("暂不支持RandomSub")
	case GossipSub:
		fallthrough
	default:
		ps, err = NewGossipSub(dsn.ctx, dsn.host, pubsubOpts...)
		if err == nil {
			logrus.Info("[DSN] gossip-sub 服务已启动")
		}
	}

	if err != nil {
		return err
	}

	dsn.pubsub = ps
	atomic.StoreInt32(&dsn.startUp, 2)
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

	logrus.Infof("[DSN] 订阅主题 [%s]", topic)

	// 如果需要订阅，则返回订阅实例
	if subscribe {
		// 设置订阅选项，包括缓冲区大小
		opts := []SubOpt{
			WithBufferSize(2048), // 设置合适的缓冲区大小
		}
		return t.Subscribe(opts...)
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

// topicSubLoop 处理订阅主题的消息循环
//
// 参数:
// - topicSub: *Subscription 表示订阅的主题
// - handler: PubSubMsgHandler 表示处理消息的回调函数
//
// 功能:
// - 持续监听订阅主题的消息
// - 过滤掉自己发送的消息
// - 异步处理接收到的消息
// - 处理各种错误情况
func (dsn *DSN) topicSubLoop(topicSub *Subscription, handler PubSubMsgHandler) {
	logrus.Info("[DSN] 开始订阅消息循环")

	for {
		// 获取下一条消息
		logrus.Info("[DSN] 等待接收下一条消息...")
		message, err := topicSub.Next(dsn.ctx)

		// 错误处理
		if err != nil {
			if err.Error() == "subscription cancelled" {
				logrus.Warn("[DSN] 订阅已被取消，退出消息循环: ", err)
				break
			}
			if err.Error() == "context canceled" {
				logrus.Info("[DSN] 上下文已取消，退出消息循环")
				break
			}
			logrus.Errorf("[DSN] 接收消息失败: %s", err.Error())
			continue
		}

		// 消息有效性检查
		if message == nil {
			logrus.Info("[DSN] 收到空消息，跳过处理")
			continue
		}

		// 忽略自己发送的消息
		if message.GetFrom() == dsn.host.ID() {
			logrus.Info("[DSN] 忽略自己发送的消息")
			continue
		}

		// 检查消息来源
		if len(message.From) == 0 {
			logrus.Info("[DSN] 消息来源为空，跳过处理")
			continue
		}

		// 解析发送者ID
		pid, err := peer.IDFromBytes(message.From)
		if err != nil {
			logrus.Errorf("[DSN] 无法解析消息发送者ID: %s", err.Error())
			continue
		}

		// 再次确认消息不是自己发送的
		if pid.String() == dsn.host.ID().String() {
			logrus.Info("[DSN] 再次确认消息来自自己，跳过处理")
			continue
		}

		// 创建处理消息的上下文
		msgCtx, cancel := context.WithTimeout(dsn.ctx, 15*time.Second)
		logrus.Infof("[DSN] 开始处理来自节点 %s 的消息", pid.String())

		// 异步处理消息
		go func(ctx context.Context, msg *Message) {
			defer cancel()
			select {
			case <-ctx.Done():
				logrus.Warn("[DSN] 消息处理超时或上下文已取消")
				return
			default:
				logrus.Info("[DSN] 调用消息处理函数")
				handler(msg)
				logrus.Info("[DSN] 消息处理完成")
			}
		}(msgCtx, message)
	}

	logrus.Info("[DSN] 订阅消息循环结束")
}

// Pubsub 返回 PubSub 实例
// 返回:
//   - *PubSub: 当前DSN例使用的PubSub实例
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
