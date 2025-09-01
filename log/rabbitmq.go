package log

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
	"github.com/chenbaoding2818/chainly/mq"
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
	ctx context.Context
	wg  *sync.WaitGroup
	// producer *RabbitMQProducer
	producer iface.IProducer
}

func NewRabbitMQ(ctx context.Context, wg *sync.WaitGroup, cfg config.OperationLog, mgCfg config.MsgQueue) *RabbitMQ {
	rmq := &RabbitMQ{
		ctx:          ctx,
		bufferChanel: make(chan []byte, cfg.BufferChanelSize),
		batchSize:    cfg.BatchSize,
		wg:           wg,
	}

	rmq.producer = mq.NewRabbitMQProducer(ctx, mgCfg.RabbitMQ)
	// 启动消费协程
	if cfg.BatchSize > 1 {
		// 1分钟超时 1分钟内如果没有累计到足够的消息，则强制发送
		rmq.ticker = time.NewTicker(flushInterval)
		rmq.runBatch()
	} else {
		rmq.run()
	}
	return rmq
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
					_msg := r.producer.CombineMessages(msgs)
					r.producer.SendAsync(r.ctx, _msg)
					msgs = msgs[:0]
				}
			case <-r.ticker.C: // 如果ticker到达时，没有累积到足够的消息，则强制发送保证一定的时效
				_msg := r.producer.CombineMessages(msgs)
				r.producer.SendAsync(r.ctx, _msg)
			case <-r.ctx.Done():
				// 关闭通道时，将剩余的消息发送出去
				close(r.bufferChanel)
				_msg := r.producer.CombineMessages(msgs)
				r.producer.SendAsync(r.ctx, _msg)
				for msg := range r.bufferChanel {
					r.producer.SendAsync(r.ctx, msg)
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
				// r.public(r.NewMQMsg(msg))
				r.producer.SendAsync(r.ctx, msg)
			case <-r.ctx.Done():
				close(r.bufferChanel)
				// 关闭通道时，将剩余的消息发送出去
				for msg := range r.bufferChanel {
					// r.public(r.NewMQMsg(msg))
					r.producer.SendAsync(r.ctx, msg)
				}
			}
		}
	}()
}
