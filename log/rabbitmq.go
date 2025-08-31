package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/chenbaoding2818/chainly/config"
	amqp "github.com/rabbitmq/amqp091-go"
	// "github.com/streadway/amqp"
)

var (
	ErrMsgQueueFull error = errors.New("msg queue full")
	flushInterval         = 1 * time.Minute // 最大等待间隔
)

type RabbitMQProducer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	reconnect chan struct{}
	urls      []string
	lock      sync.Mutex
}

func newRabbitMQProducer(cfg *config.RabbitMQConfig) *RabbitMQProducer {
	return &RabbitMQProducer{
		urls:      cfg.Urls,
		reconnect: make(chan struct{}),
	}
}

type RabbitMQ struct {
	bufferChanel chan []byte
	batchSize    int32
	ticker       *time.Ticker
	// 通道关闭时 发送剩余的消息
	ctx      context.Context
	wg       *sync.WaitGroup
	producer *RabbitMQProducer
}

func NewRabbitMQ(ctx context.Context, wg *sync.WaitGroup, cfg config.OperationLog, mgCfg config.MsgQueue) *RabbitMQ {
	mq := &RabbitMQ{
		ctx:          ctx,
		bufferChanel: make(chan []byte, cfg.BufferChanelSize),
		batchSize:    cfg.BatchSize,
		wg:           wg,
	}

	mq.producer = newRabbitMQProducer(mgCfg.RabbitMQ)
	// 连接rabbitMQ
	go mq.connect()
	// 重连
	go mq.reconnect()
	// 启动消费协程
	if cfg.BatchSize > 1 {
		// 1分钟超时 1分钟内如果没有累计到足够的消息，则强制发送
		mq.ticker = time.NewTicker(flushInterval)
		mq.runBatch()
	} else {
		mq.run()
	}
	return mq
}

func (r *RabbitMQ) SendMsg(msg []byte) error {
	select {
	case r.bufferChanel <- msg:
		return nil
	default: // 消息累积满后直接丢弃
		// TODO: 缓冲满的时候可以做一些逻辑，比如丢弃一些消息，或者等待一段时间再发送
		return ErrMsgQueueFull
	}
}

// run 运行消费协程
func (r *RabbitMQ) runBatch() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		msgs := make([][]byte, 0, r.batchSize)
		for {
			select {
			case msg := <-r.bufferChanel:
				// 批量处理
				if len(msgs) < int(r.batchSize) {
					msgs = append(msgs, msg)
				} else {
					_msg := r.combineMessages(msgs)
					r.public(r.NewMQMsg(_msg))
					msgs = msgs[:0]
				}
			case <-r.ticker.C: // 如果ticker到达时，没有累积到足够的消息，则强制发送保证一定的时效
				r.public(r.NewMQMsg(r.combineMessages(msgs)))
			case <-r.ctx.Done():
				// 关闭通道时，将剩余的消息发送出去
				close(r.bufferChanel)
				r.public(r.NewMQMsg(r.combineMessages(msgs)))
				for msg := range r.bufferChanel {
					r.public(r.NewMQMsg(msg))
				}
				return
			}
		}
	}()
}

// run 运行消费协程
func (r *RabbitMQ) run() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case msg := <-r.bufferChanel:
				r.public(r.NewMQMsg(msg))
			case <-r.ctx.Done():
				close(r.bufferChanel)
				// 关闭通道时，将剩余的消息发送出去
				for msg := range r.bufferChanel {
					r.public(r.NewMQMsg(msg))
				}
			}
		}
	}()
}

// combineMessages 合并多条消息为批量格式
func (r *RabbitMQ) combineMessages(msgs [][]byte) []byte {
	var buf bytes.Buffer
	for _, msg := range msgs {
		buf.Write(msg)
	}
	return buf.Bytes()
}

func (r *RabbitMQ) public(body amqp.Publishing) {
	// 判断body是否为空
	if len(body.Body) == 0 {
		return
	}
	r.producer.lock.Lock()
	defer r.producer.lock.Unlock()
}

func (r *RabbitMQ) NewMQMsg(body []byte) amqp.Publishing {
	// 空数据检测
	if len(body) == 0 {
		return amqp.Publishing{}
	}
	return amqp.Publishing{
		ContentType:  "text/plain",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		MessageId:    fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
}

func (r *RabbitMQ) connect() {
	r.producer.lock.Lock()
	defer r.producer.lock.Unlock()
	for _, url := range r.producer.urls {
		conn, err := amqp.Dial(url)
		if err == nil {
			r.producer.conn = conn
			ch, err := conn.Channel()
			if err == nil {
				r.producer.channel = ch
				// 声明交换机
				err = ch.ExchangeDeclare(
					"logs",
					"fanout",
					true,
					false,
					false,
					false,
					nil,
				)
				if err == nil {

				}

			}

		}
	}

	// 重连
	r.producer.reconnect <- struct{}{}
}

func (r *RabbitMQ) reconnect() error {
	for {
		select {
		case <-r.producer.reconnect:
			const maxRetries = 5
			for i := 0; i < maxRetries; i++ {
				wait := time.Duration(i*i) * time.Second
				time.Sleep(wait)
				r.connect()
				if r.producer.conn != nil {
					break
				} else if i == maxRetries-1 {

				}
			}

			r.connect()
		case <-r.ctx.Done():
			return nil
		}
	}
}
