package mq

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PendingMessage struct {
	DeliveryTag     uint64
	msg             []byte
	sendAt          time.Time
	successCallback iface.WaitForConfirmFunc
	failureCallback iface.WaitForConfirmFunc
}

type RabbitMQProducer struct {
	cfg             *config.RabbitMQConfig
	conn            *amqp.Connection
	channel         *amqp.Channel
	connCloseChan   chan *amqp.Error
	connErrChan     chan struct{}
	ctx             context.Context // 退出信号上下文
	wg              *sync.WaitGroup // 框架层退出等待组
	_wg             sync.WaitGroup
	config          amqp.Config               // rabbitmq配置
	confirmations   chan amqp.Confirmation    // 确认通道
	pendingLock     sync.Mutex                // 保护pending消息的锁
	pendingMessages map[uint64]PendingMessage // 待确认消息
	msgSequence     uint64                    // 消息序列号
}

func NewRabbitMQProducer(ctx context.Context, wg *sync.WaitGroup, cfg *config.RabbitMQConfig) iface.IProducer {
	rbmq := &RabbitMQProducer{
		cfg: cfg,
		config: amqp.Config{
			Heartbeat: time.Duration(cfg.HeartBeat) * time.Second,
			Dial:      amqp.DefaultDial(5 * time.Second),
		},
		ctx:             ctx,
		wg:              wg,
		pendingMessages: make(map[uint64]PendingMessage),
	}
	// 创建连接
	if err := rbmq.Connect(); err != nil {
		panic(err)
	}
	// 启动重连机制
	rbmq.wg.Add(1)
	go rbmq.ReconnectMonitor()
	return rbmq
}

func (r *RabbitMQProducer) Connect() error {
	// 连接集群 目前负载均衡采用顺序连接 配置写死 TODO:后期可加上服务发现机制获取集群地址
	var err error
	for _, url := range r.cfg.Urls {
		r.conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
	}
	if err != nil {
		return errors.New("connect rabbitmq failed")
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		r.conn.Close()
		return fmt.Errorf("channel open failed: %s", err.Error())
	}
	// 开启confirm模式
	if r.cfg.ConfirmEnable {
		if err := r.channel.Confirm(false); err != nil {
			r.channel.Close()
			r.conn.Close()
			return fmt.Errorf("confirm mode enable failed: %w", err)
		}
		r.confirmations = r.channel.NotifyPublish(make(chan amqp.Confirmation, r.cfg.MaxPendingMessage))
		r._wg.Add(2)
		// 开启确认协程
		go r.waitForConfirmation()
		// 开启超时监控协程
		go r.pendingMsgTimeoutMonitor()
	}

	if err := r.channel.ExchangeDeclare(
		r.cfg.ExchangeName,  // 交换机名称
		r.cfg.ExchangeType,  // 交换机类型
		r.cfg.ConfirmEnable, // 持久化
		false,               // 自动删除
		false,               // internal
		true,                // no-wait
		nil,                 // table
	); err != nil {
		r.Close()
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	// 监听连接关闭
	r.connCloseChan = r.conn.NotifyClose(make(chan *amqp.Error))
	return nil
}

// Reconnect 重连
func (r *RabbitMQProducer) Reconnect() error {
	// 采取无限重连的方式，直到重连成功
	for {
		select {
		case <-r.ctx.Done():
			return nil
		default:
			if err := r.Connect(); err == nil {
				return err
			}

			select {
			case <-time.After(time.Second * time.Duration(r.cfg.HeartBeat)):
				continue
			}
		}
	}
}

// ReconnectMonitor 重连监控
func (r *RabbitMQProducer) ReconnectMonitor() error {
	defer r.wg.Done()
	for {
		select {
		case <-r.connCloseChan: // 连接关闭时
			// 开始进行重连时先关闭当前连接
			r.Close()
			r.Reconnect()
		case <-r.ctx.Done(): // 重连退出
			r.Close()
			return nil
		}
	}
}

// waitForConfirmation 等待确认
func (r *RabbitMQProducer) waitForConfirmation() {
	defer r._wg.Done()
	for {
		select {
		case confirm, ok := <-r.confirmations:
			if !ok {
				return
			}
			r.pendingLock.Lock()
			pendingMsg, ok := r.pendingMessages[confirm.DeliveryTag]
			if ok {
				delete(r.pendingMessages, confirm.DeliveryTag)
			}
			r.pendingLock.Unlock()
			if !ok {
				continue
			}
			// mq处理成功回调
			if confirm.Ack && pendingMsg.successCallback != nil {
				// 处理消息的成功回调
				if pendingMsg.successCallback != nil {
					pendingMsg.successCallback()
				}
			} else { // mq处理失败回调
				if pendingMsg.failureCallback != nil {
					pendingMsg.failureCallback()
				}
			}
		case <-r.connCloseChan:
			return
		}
	}
}

// pendingMsgTimeoutMonitor 监控待确认消息超时
func (r *RabbitMQProducer) pendingMsgTimeoutMonitor() {
	defer r._wg.Done()
	for {
		select {
		case <-time.After(time.Second * 5):
			for seq, msg := range r.pendingMessages {
				// 如果超过超时时间，则认为消息超时，调用失败回调
				if time.Since(msg.sendAt) > time.Duration(r.cfg.MsgConfirmTimeout) {
					if msg.failureCallback != nil {
						msg.failureCallback()
					}
					r.pendingLock.Lock()
					delete(r.pendingMessages, seq)
					r.pendingLock.Unlock()
				}
			}
		case <-r.connCloseChan:
			return

		}
	}
}

func (r *RabbitMQProducer) Close() {
	close(r.connCloseChan)
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	// 等待内部协程退出
	r._wg.Wait()
}

func (r *RabbitMQProducer) SendWithConfirm(dst string, msg []byte, successCallback, failureCallback iface.WaitForConfirmFunc) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	// 空数据检测
	if len(msg) == 0 {
		return errors.New("empty message")
	}
	// 检测是否开启confirm模式
	if !r.cfg.ConfirmEnable {
		r.SendWithoutConfirm(dst, msg)
		return nil
	}
	seq := atomic.AddUint64(&r.msgSequence, 1)
	publishing := r.NewMQMsg(msg, seq)
	err := r.channel.Publish(
		r.cfg.ExchangeName, // 交换机名称
		dst,                // 路由键
		false,              // mandatory
		false,              // immediate
		publishing,         // msg
	)
	if err != nil {
		return err
	}
	// 将消息放入待确认消息队列
	r.pendingLock.Lock()
	pendingMsg := PendingMessage{
		DeliveryTag:     seq,
		msg:             msg,
		sendAt:          time.Now(),
		successCallback: successCallback,
		failureCallback: failureCallback,
	}
	r.pendingMessages[seq] = pendingMsg
	r.pendingLock.Unlock()

	return nil
}

// SendAsync 不需要等待确认
func (r *RabbitMQProducer) SendWithoutConfirm(dst string, msg []byte) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	// 空数据检测
	if len(msg) == 0 {
		return errors.New("empty message")
	}
	// 检测是否开启confirm模式 如果是confirm模式，则使用SendWithConfirm
	if r.cfg.ConfirmEnable {
		r.SendWithConfirm(dst, msg, nil, nil)
		return nil
	}
	publishing := r.NewMQMsg(msg, uint64(time.Now().UnixMilli()))
	err := r.channel.Publish(
		r.cfg.ExchangeName, // 交换机名称
		dst,                // 路由键
		false,              // mandatory
		false,              // immediate
		publishing,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQProducer) CombineMessages(msgs [][]byte) []byte {
	var buf bytes.Buffer
	for _, msg := range msgs {
		buf.Write(msg)
	}
	return buf.Bytes()
}

func (r *RabbitMQProducer) NewMQMsg(body []byte, seq uint64) amqp.Publishing {
	return amqp.Publishing{
		ContentType:  "text/plain",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		MessageId:    fmt.Sprintf("%d", seq),
	}
}
