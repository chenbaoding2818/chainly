package mq

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	cfg            *config.RabbitMQConfig
	conn           *amqp.Connection
	channel        *amqp.Channel
	lock           sync.Mutex
	ctx            context.Context
	config         amqp.Config            // rabbitmq配置
	confirmEnabled bool                   // 是否开启confirm模式
	confirmations  chan amqp.Confirmation // 确认通道
	reconnect      chan struct{}          // 重连信号
	reconnectCount int                    // 重连次数
	reconnectMax   int                    // 最大重连次数
}

func NewRabbitMQProducer(ctx context.Context, rabbitMQCfg *config.RabbitMQConfig) iface.IProducer {
	return &RabbitMQProducer{
		cfg:          rabbitMQCfg,
		ctx:          ctx,
		reconnect:    make(chan struct{}),
		reconnectMax: 10,
	}
}

func (r *RabbitMQProducer) Connect(exchangeName string, exchangeType string) error {
	// 连接集群 目前负载均衡采用顺序连接 配置写死 TODO:后期可加上服务发现机制获取集群地址
	r.lock.Lock()
	defer r.lock.Unlock()
	// 可用连接
	var (
		err       error
		vaildConn bool = false
	)

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
	if r.confirmEnabled {
		r.confirmations = make(chan amqp.Confirmation, 100)
	}

	if err := r.channel.ExchangeDeclare(
		exchangeName, // 交换机名称
		exchangeType, // 交换机类型
		true,         // 持久化
		false,        // 自动删除
		false,        // internal
		true,         // no-wait
		nil,          // table
	); err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	// 连接关闭
	connClose := r.conn.NotifyClose(make(chan *amqp.Error))
	select {
	case err := <-connClose:
		if err != nil {
			// log.Errorf("[error] %s connection closed: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		}
		// 进入重连流程
		r.reconnect <- struct{}{}
		return err
	case <-r.ctx.Done():
		return nil
	}
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

func (r *RabbitMQProducer) Send(ctx context.Context, msg []byte, f iface.WaitForConfirmFunc) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	publishing := r.NewMQMsg(msg)
	fmt.Printf("publishing: %v\n", publishing)
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

	err = r.waitForConfirmation(0, f)
	if err != nil {
		return err
	}

	return nil
}

// SendAsync 不需要等待确认
func (r *RabbitMQProducer) SendAsync(ctx context.Context, msg []byte) error {
	// 空数据检测
	if len(msg) == 0 {
		return errors.New("empty message")
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

func (r *RabbitMQProducer) waitForConfirmation(seq uint64, f iface.WaitForConfirmFunc) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	for {
		select {
		case confirm := <-r.confirmations:
			if confirm.DeliveryTag == seq {
				if !confirm.Ack {
					return fmt.Errorf("message %d nack-ed by broker", seq)
				}
				return nil
			}
		case <-ctx.Done():
			if f != nil {
				// 执行回调函数
				return f()
			}
			return fmt.Errorf("confirmation timeout for message %d", seq)
		}
	}
}

func (r *RabbitMQProducer) Close() {
}
