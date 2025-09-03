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
	DeliveryTag     int64
	sendAt          time.Time
	successCallback iface.WaitForConfirmFunc
	failureCallback iface.WaitForConfirmFunc
}

type RabbitMQProducer struct {
	cfg             *config.RabbitMQConfig
	conn            *amqp.Connection
	channel         *amqp.Channel
	connCloseChan   chan *amqp.Error
	connDoneChan    chan struct{}
	lock            sync.Mutex
	ctx             context.Context
	config          amqp.Config              // rabbitmq配置
	confirmations   chan amqp.Confirmation   // 确认通道
	reconnect       chan struct{}            // 重连信号
	pendingLock     sync.Mutex               // 保护pending消息的锁
	pendingMessages map[int64]PendingMessage // 待确认消息
	msgSequence     int64                    // 消息序列号
}

func NewRabbitMQProducer(ctx context.Context, cfg *config.RabbitMQConfig) iface.IProducer {
	rbmq := &RabbitMQProducer{
		cfg: cfg,
		config: amqp.Config{
			Heartbeat: time.Duration(cfg.HeartBeat) * time.Second,
			Dial:      amqp.DefaultDial(5 * time.Second),
		},
		ctx:             ctx,
		reconnect:       make(chan struct{}),
		pendingMessages: make(map[int64]PendingMessage),
	}

	if 

	return rbmq
}

func (r *RabbitMQProducer) Connect() error {
	// 连接集群 目前负载均衡采用顺序连接 配置写死 TODO:后期可加上服务发现机制获取集群地址
	// 可用连接
	var (
		err       error
		vaildConn bool = false
	)
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, url := range r.cfg.Urls {
		r.conn, err = amqp.Dial(url)
		if err == nil {
			vaildConn = true
			break
		}
	}
	if !vaildConn {
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
			return fmt.Errorf("confirm mode enable failed: %w", err)
		}
		r.confirmations = r.channel.NotifyPublish(make(chan amqp.Confirmation, r.cfg.MaxPendingMessage))
		go r.waitForConfirmation()
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
		r.conn.Close()
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	// 监听连接关闭
	r.connCloseChan = r.conn.NotifyClose(make(chan *amqp.Error))
	return nil
}

func (r *RabbitMQProducer) Reconnect() error {
	for {
		select {
		case <-r.reconnect:
			return nil
		case <-r.ctx.Done(): // 重连退出
			return nil
		default:
		}
	}
}

func (r *RabbitMQProducer) SendWithConfirm(dst string, msg []byte, successCallback, failureCallback iface.WaitForConfirmFunc) error {
	// 空数据检测
	if len(msg) == 0 {
		return errors.New("empty message")
	}
	// 检测是否开启confirm模式
	if !r.cfg.ConfirmEnable {
		r.SendWithoutConfirm(dst, msg)
		return nil
	}
	r.lock.Lock()
	publishing := r.NewMQMsg(msg)
	fmt.Printf("publishing: %v\n", publishing)
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
	r.lock.Unlock()

	// 将消息放入待确认消息队列
	r.pendingLock.Lock()
	seq := atomic.AddInt64(&r.msgSequence, 1)
	pendingMsg := PendingMessage{
		DeliveryTag:     seq,
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
	// 空数据检测
	if len(msg) == 0 {
		return errors.New("empty message")
	}
	// 检测是否开启confirm模式 如果是confirm模式，则使用SendWithConfirm
	if r.cfg.ConfirmEnable {
		r.SendWithConfirm(dst, msg, nil, nil)
		return nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	publishing := r.NewMQMsg(msg)
	err := r.channel.Publish(
		"",    // 交换机名称
		"",    // 路由键
		false, // mandatory
		false, // immediate
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

func (r *RabbitMQProducer) NewMQMsg(body []byte) amqp.Publishing {
	return amqp.Publishing{
		ContentType:  "text/plain",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		MessageId:    fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
}

// waitForConfirmation 等待确认
func (r *RabbitMQProducer) waitForConfirmation() {
	for {
		select {
		case confirm := <-r.confirmations:
			// if confirm.DeliveryTag == seq {
			// 	if !confirm.Ack {
			// 		return fmt.Errorf("message %d nack-ed by broker", seq)
			// 	}
			// 	return nil
			// }
			fmt.Println("confirm:", confirm)
		case <-r.connDoneChan: // 退出等待
			return
		case <-r.ctx.Done(): // 退出等待
			return
		}
	}
}

// pendingMsgTimeoutMonitor 监控待确认消息超时
func (r *RabbitMQProducer) pendingMsgTimeoutMonitor() {
	for {
		select {
		case <-time.After(time.Second * 1):
			r.pendingLock.Lock()
			for seq, msg := range r.pendingMessages {
				// 如果超过超时时间，则认为消息超时，调用失败回调
				if time.Since(msg.sendAt) > time.Duration(r.cfg.MsgConfirmTimeout) {
					fmt.Println("pending message timeout:", seq)
					if msg.failureCallback != nil {
						msg.failureCallback()
					}
					delete(r.pendingMessages, seq)
				}
			}
			r.pendingLock.Unlock()
		case <-r.connDoneChan: // 退出监控
			return
		case <-r.ctx.Done(): // 退出监控
			return
			// case
		}
	}
}

func (r *RabbitMQProducer) Close() {
	close(r.connDoneChan)
}
