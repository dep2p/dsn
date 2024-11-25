// package dsn 提供了发布订阅功能的实现
package dsn

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/sirupsen/logrus"
)

// NodeStatus 表示节点的状态
type NodeStatus int

const (
	Online     NodeStatus = iota // 节点在线
	Suspicious                   // 节点可疑
	Offline                      // 节点离线
)

// StatusChange 表示节点状态的变化
type StatusChange struct {
	PeerID            peer.ID       // 节点ID
	OldStatus         NodeStatus    // 旧状态
	NewStatus         NodeStatus    // 新状态
	Timestamp         time.Time     // 状态变化时间戳
	Score             float64       // 节点评分
	ConnectionQuality float64       // 连接质量
	CheckInterval     time.Duration // 检查间隔
	FailedAttempts    int           // 连续失败尝试次数
	NetworkLatency    time.Duration // 网络延迟
}

// Node 表示一个节点及其状态信息
type Node struct {
	ID                  peer.ID         // 节点ID
	Status              NodeStatus      // 当前状态
	LastSeen            time.Time       // 最后一次看到节点的时间
	FailedAttempts      int             // 连续失败尝试次数
	Score               float64         // 节点评分
	History             []StatusChange  // 状态变化历史
	CheckInterval       time.Duration   // 检查间隔
	SuccessiveSuccesses int             // 连续成功次数
	LastCheckTime       time.Time       // 上次检查时间
	ConnectionQuality   float64         // 连接质量
	NetworkCondition    []time.Duration // 网络延迟历史
}

// NodeStatusTracker 用于跟踪和管理节点状态
type NodeStatusTracker struct {
	host                 host.Host           // libp2p主机
	nodes                map[peer.ID]*Node   // 节点映射
	defaultCheckInterval time.Duration       // 默认检查间隔
	offlineThreshold     int                 // 离线阈值
	pingTimeout          time.Duration       // ping超时时间
	pingService          *ping.PingService   // ping服务
	mu                   sync.RWMutex        // 读写锁
	ctx                  context.Context     // 上下文
	cancel               context.CancelFunc  // 取消函数
	statusChanges        chan StatusChange   // 状态变化通道
	subscribers          []chan StatusChange // 订阅者列表
}

// NewNodeStatusTracker 创建一个新的 NodeStatusTracker 实例
// 参数:
//   - h: libp2p主机
//   - defaultCheckInterval: 默认检查间隔
//   - offlineThreshold: 离线阈值
//   - pingTimeout: ping超时时间
//
// 返回:
//   - *NodeStatusTracker: 新创建的NodeStatusTracker实例
func NewNodeStatusTracker(h host.Host, defaultCheckInterval time.Duration, offlineThreshold int, pingTimeout time.Duration) *NodeStatusTracker {
	// 创建一个带有取消功能的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 返回一个新的NodeStatusTracker实例
	return &NodeStatusTracker{
		host:                 h,
		nodes:                make(map[peer.ID]*Node),
		defaultCheckInterval: defaultCheckInterval,
		offlineThreshold:     offlineThreshold,
		pingTimeout:          pingTimeout,
		pingService:          ping.NewPingService(h),
		ctx:                  ctx,
		cancel:               cancel,
		statusChanges:        make(chan StatusChange, 100),
	}
}

// Start 启动节点状态跟踪器
func (nst *NodeStatusTracker) Start() {
	go nst.checkNodesRoutine()    // 启动节点检查例程
	go nst.processStatusChanges() // 启动状态变化处理例程
}

// Stop 停止节点状态跟踪器
func (nst *NodeStatusTracker) Stop() {
	nst.cancel() // 调用取消函数停止所有goroutine
}

// checkNodesRoutine 定期检查所有节点的状态
func (nst *NodeStatusTracker) checkNodesRoutine() {
	for {
		select {
		case <-nst.ctx.Done():
			return // 如果上下文被取消，退出循环
		default:
			nst.checkAllNodes()                  // 检查所有节点
			time.Sleep(nst.defaultCheckInterval) // 等待默认检查间隔
		}
	}
}

// checkAllNodes 检查所有已知节点的状态
func (nst *NodeStatusTracker) checkAllNodes() {
	// 获取读锁以安全地读取节点列表
	nst.mu.RLock()
	peers := make([]peer.ID, 0, len(nst.nodes))
	for pid := range nst.nodes {
		peers = append(peers, pid)
	}
	nst.mu.RUnlock()

	// 遍历所有节点并检查它们的状态
	for _, pid := range peers {
		nst.checkNode(pid) // 检查每个节点
	}
}

// checkNode 检查单个节点的状态
// 参数:
//   - pid: 要检查的节点ID
func (nst *NodeStatusTracker) checkNode(pid peer.ID) {
	nst.mu.Lock()
	defer nst.mu.Unlock()

	// 基础配置 - 更友好的参数设置
	const (
		maxScore = 1.0 // 最大评分
		minScore = 0.3 // 最小评分不要太低，给予恢复机会
		// 成功恢复的奖励大于失败的惩罚
		scoreIncrement = 0.2 // 成功时的评分增加
		scoreDecrement = 0.1 // 失败时的评分减少
		// 降低恢复门槛
		recoveryThreshold = 2 // 恢复所需的连续成功次数
		// 网络波动容忍度
		networkJitterTolerance = 3 * time.Second
	)

	// 获取或创建节点
	node, exists := nst.nodes[pid]
	if !exists {
		node = &Node{
			ID:                  pid,
			Status:              Online,
			LastSeen:            time.Now(),
			Score:               0.7, // 初始分值适中，给予成长空间
			CheckInterval:       nst.defaultCheckInterval,
			FailedAttempts:      0,
			History:             make([]StatusChange, 0),
			SuccessiveSuccesses: 0,
			LastCheckTime:       time.Now(),
			ConnectionQuality:   0.7,                      // 初始连接质量适中
			NetworkCondition:    make([]time.Duration, 0), // 存储最近的网络延迟
		}
		nst.nodes[pid] = node
	}

	oldStatus := node.Status
	checkTime := time.Now()
	// timeSinceLastCheck := checkTime.Sub(node.LastCheckTime)

	// 检查网络状态
	isConnected := nst.host.Network().Connectedness(pid) == network.Connected
	var pingLatency time.Duration

	// 智能 ping 策略
	if !isConnected {
		// 如果距离上次成功不久，给予宽限
		if time.Since(node.LastSeen) < networkJitterTolerance {
			// 临时网络抖动，保持当前状态
			return
		}

		start := time.Now()
		canPing := nst.pingNode(pid)
		if canPing {
			pingLatency = time.Since(start)
			isConnected = true // 如果可以 ping 通，视为连接正常
		}
	}

	// 状态更新逻辑
	if isConnected {
		// 成功连接处理
		node.SuccessiveSuccesses++
		node.LastSeen = checkTime

		// 智能失败次数重置
		if node.FailedAttempts > 0 {
			// 渐进式减少失败计数，而不是直接清零
			node.FailedAttempts = int(float64(node.FailedAttempts) * 0.5)
		}

		// 更新网络质量评估
		if pingLatency > 0 {
			node.NetworkCondition = append(node.NetworkCondition, pingLatency)
			if len(node.NetworkCondition) > 10 {
				node.NetworkCondition = node.NetworkCondition[1:]
			}

			// 计算平均网络状况
			var avgLatency time.Duration
			for _, latency := range node.NetworkCondition {
				avgLatency += latency
			}
			avgLatency /= time.Duration(len(node.NetworkCondition))

			// 基于平均延迟动态调整连接质量
			qualityFactor := float64(time.Second) / float64(avgLatency)
			node.ConnectionQuality = (node.ConnectionQuality*0.7 + qualityFactor*0.3)
		}

		// 评分提升 - 根据连续成功次数给予额外奖励
		scoreBonus := scoreIncrement * (1 + float64(node.SuccessiveSuccesses)/10)
		node.Score = math.Min(maxScore, node.Score+scoreBonus)

		// 状态恢复逻辑 - 更宽容的恢复策略
		if node.Status != Online && node.SuccessiveSuccesses >= recoveryThreshold {
			nst.updateNodeStatus(node, Online)
		}
	} else {
		// 失败处理 - 更友好的策略
		node.FailedAttempts++
		node.SuccessiveSuccesses = 0

		// 评分调整 - 考虑历史表现
		penaltyFactor := math.Max(0.5, float64(node.Score)) // 评分越高，惩罚越轻
		actualDecrement := scoreDecrement * penaltyFactor
		node.Score = math.Max(minScore, node.Score-actualDecrement)

		// 状态降级逻辑 - 更渐进的方式
		if node.FailedAttempts >= nst.offlineThreshold*3 {
			if node.Status != Offline {
				nst.updateNodeStatus(node, Offline)
				// 异步进行深度检查
				go nst.performDeepCheck(pid)
			}
		} else if node.FailedAttempts >= nst.offlineThreshold {
			if node.Status != Suspicious {
				nst.updateNodeStatus(node, Suspicious)
			}
		}
	}

	// 动态调整检查间隔 - 基于网络状况
	nst.adjustCheckIntervalBasedOnNetwork(node)

	// 记录状态变化
	if oldStatus != node.Status {
		statusChange := StatusChange{
			PeerID:            pid,
			OldStatus:         oldStatus,
			NewStatus:         node.Status,
			Timestamp:         checkTime,
			Score:             node.Score,
			ConnectionQuality: node.ConnectionQuality,
			CheckInterval:     node.CheckInterval,
			FailedAttempts:    node.FailedAttempts,
			NetworkLatency:    pingLatency,
		}

		node.History = append(node.History, statusChange)
		if len(node.History) > 100 {
			node.History = node.History[1:]
		}

		nst.recordStatusChange(statusChange)
	}

	// 更新检查时间
	node.LastCheckTime = checkTime

	// 调试日志
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.Debugf("节点检查 [%s] - 状态: %v, 评分: %.2f, 质量: %.2f, 失败: %d, 连续成功: %d",
			pid.String()[:8], node.Status, node.Score, node.ConnectionQuality,
			node.FailedAttempts, node.SuccessiveSuccesses)
	}
}

// updateNodeStatus 更新节点状态
// 参数:
//   - node: 要更新的节点
//   - status: 新的状态
func (nst *NodeStatusTracker) updateNodeStatus(node *Node, status NodeStatus) {
	node.Status = status
	if status == Online {
		node.LastSeen = time.Now()
		node.FailedAttempts = 0
		node.Score = min(node.Score+0.1, 1.0) // 增加节点评分，最高为1.0
	} else {
		node.Score = max(node.Score-0.2, 0.0) // 减少节点评分，最低为0.0
	}
}

// pingNode 尝试 ping 指定的节点
// 参数:
//   - pid: 要ping的节点ID
//
// 返回:
//   - bool: ping是否成功
func (nst *NodeStatusTracker) pingNode(pid peer.ID) bool {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(nst.ctx, nst.pingTimeout)
	defer cancel()

	// 执行ping操作
	result := nst.pingService.Ping(ctx, pid)
	select {
	case <-result:
		return true // ping成功
	case <-ctx.Done():
		return false // ping超时
	}
}

// adjustCheckInterval 根据节点历史表现调整检查间隔
// 参数:
//   - node: 要调整的节点
func (nst *NodeStatusTracker) adjustCheckInterval(node *Node) {
	if node.Score > 0.8 {
		node.CheckInterval = nst.defaultCheckInterval * 2 // 评分高，增加检查间隔
	} else if node.Score < 0.3 {
		node.CheckInterval = nst.defaultCheckInterval / 2 // 评分低，减少检查间隔
	} else {
		node.CheckInterval = nst.defaultCheckInterval // 评分中等，使用默认间隔
	}
}

// recordStatusChange 记录状态变化并通知订阅者
// 参数:
//   - change: 状态变化信息
func (nst *NodeStatusTracker) recordStatusChange(change StatusChange) {
	select {
	case nst.statusChanges <- change:
		// 成功将变化发送到通道
	default:
		// 如果通道已满，记录日志但不阻塞
		log.Warnf("状态变化通道已满，丢弃变化：%v", change)
	}
}

// processStatusChanges 处理状态变化并通知订阅者
func (nst *NodeStatusTracker) processStatusChanges() {
	for {
		select {
		case change := <-nst.statusChanges:
			nst.notifySubscribers(change) // 通知所有订阅者
		case <-nst.ctx.Done():
			return // 如果上下文被取消，退出循环
		}
	}
}

// notifySubscribers 通知所有订阅者状态变化
// 参数:
//   - change: 状态变化信息
func (nst *NodeStatusTracker) notifySubscribers(change StatusChange) {
	for _, sub := range nst.subscribers {
		select {
		case sub <- change:
			// 成功发送状态变化到订阅者
		default:
			// 如果订阅者的通道已满，跳过但不阻塞
		}
	}
}

// SubscribeStatusChanges 订阅状态变化事件
// 返回:
//   - chan StatusChange: 状态变化通道
func (nst *NodeStatusTracker) SubscribeStatusChanges() chan StatusChange {
	ch := make(chan StatusChange, 10)
	nst.mu.Lock()
	nst.subscribers = append(nst.subscribers, ch)
	nst.mu.Unlock()
	return ch
}

// UnsubscribeStatusChanges 取消订阅状态变化事件
// 参数:
//   - ch: 要取消订阅的通道
func (nst *NodeStatusTracker) UnsubscribeStatusChanges(ch chan StatusChange) {
	nst.mu.Lock()
	for i, sub := range nst.subscribers {
		if sub == ch {
			// 从订阅者列表中移除指定的通道
			nst.subscribers = append(nst.subscribers[:i], nst.subscribers[i+1:]...)
			break
		}
	}
	nst.mu.Unlock()
}

// UpdateNodeStatus 允许外部组件更新节点状态
// 参数:
//   - pid: 要更新的节点ID
//   - status: 新的状态
func (nst *NodeStatusTracker) UpdateNodeStatus(pid peer.ID, status NodeStatus) {
	nst.mu.Lock()
	defer nst.mu.Unlock()

	// 获取或创建节点
	node, exists := nst.nodes[pid]
	if !exists {
		node = &Node{ID: pid, CheckInterval: nst.defaultCheckInterval}
		nst.nodes[pid] = node
	}

	oldStatus := node.Status
	nst.updateNodeStatus(node, status)

	// 如果状态发生变化，记录变化
	if oldStatus != status {
		nst.recordStatusChange(StatusChange{
			PeerID:    pid,
			OldStatus: oldStatus,
			NewStatus: status,
			Timestamp: time.Now(),
		})
	}
}

// min 返回两个浮点数中的较小值
// 参数:
//   - a, b: 要比较的两个浮点数
//
// 返回:
//   - float64: 较小的值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// max 返回两个浮点数中的较大值
// 参数:
//   - a, b: 要比较的两个浮点数
//
// 返回:
//   - float64: 较大的值
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// confirmNodeStatus 对可能不健康的节点进行额外确认
// 参数:
//   - pid: 要确认状态的节点ID
func (nst *NodeStatusTracker) confirmNodeStatus(pid peer.ID) {
	retries := 3
	retryInterval := time.Second * 30

	for i := 0; i < retries; i++ {
		time.Sleep(retryInterval)

		nst.mu.Lock()
		node, exists := nst.nodes[pid]
		if !exists {
			nst.mu.Unlock()
			return
		}

		// 检查节点是否已连接或可以ping通
		if nst.host.Network().Connectedness(pid) == network.Connected || nst.pingNode(pid) {
			nst.updateNodeStatus(node, Online)
			nst.mu.Unlock()
			return
		}

		nst.mu.Unlock()
	}

	// 如果多次重试后仍然失败，则将节点标记为离线
	nst.mu.Lock()
	defer nst.mu.Unlock()

	node, exists := nst.nodes[pid]
	if exists {
		nst.updateNodeStatus(node, Offline)
		nst.recordStatusChange(StatusChange{
			PeerID:    pid,
			OldStatus: node.Status,
			NewStatus: Offline,
			Timestamp: time.Now(),
		})
	}
}

// adjustCheckIntervalBasedOnNetwork 基于网络状况调整检查间隔
func (nst *NodeStatusTracker) adjustCheckIntervalBasedOnNetwork(node *Node) {
	// 基础间隔
	baseInterval := nst.defaultCheckInterval

	// 根据节点状态调整
	switch node.Status {
	case Online:
		// 在线状态下，根据连接质量调整
		if node.ConnectionQuality > 0.8 {
			// 连接质量好，可以适当延长间隔
			node.CheckInterval = baseInterval * 2
		} else {
			// 保持默认间隔
			node.CheckInterval = baseInterval
		}
	case Suspicious:
		// 可疑状态下，缩短检查间隔但不要太频繁
		node.CheckInterval = baseInterval / 2
	case Offline:
		// 离线状态下，采用渐进式增加的检查间隔
		// 避免频繁检查已知离线的节点
		offlineTime := time.Since(node.LastSeen)
		if offlineTime > time.Hour {
			node.CheckInterval = baseInterval * 4
		} else if offlineTime > time.Minute*30 {
			node.CheckInterval = baseInterval * 2
		} else {
			node.CheckInterval = baseInterval
		}
	}

	// 确保检查间隔在合理范围内
	minInterval := time.Second * 5
	maxInterval := time.Minute * 5
	node.CheckInterval = time.Duration(
		math.Max(
			float64(minInterval),
			math.Min(float64(maxInterval), float64(node.CheckInterval)),
		),
	)
}

// performDeepCheck 执行深度节点检查
func (nst *NodeStatusTracker) performDeepCheck(pid peer.ID) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// 尝试多种方式检查节点
	checks := []struct {
		name  string
		check func() bool
	}{
		{"直接连接", func() bool {
			return nst.host.Network().Connectedness(pid) == network.Connected
		}},
		{"Ping检查", func() bool {
			return nst.pingNode(pid)
		}},
		{"DHT查找", func() bool {
			// 如果实现了DHT，可以尝试通过DHT查找节点
			return false // 待实现
		}},
	}

	for _, check := range checks {
		select {
		case <-ctx.Done():
			return
		default:
			if check.check() {
				// 如果任何检查成功，更新节点状态
				nst.mu.Lock()
				if node, exists := nst.nodes[pid]; exists {
					nst.updateNodeStatus(node, Suspicious)         // 降级到可疑而不是直接在线
					node.FailedAttempts = nst.offlineThreshold - 1 // 给予恢复机会
				}
				nst.mu.Unlock()
				return
			}
		}
	}
}
