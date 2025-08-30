package log

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/chenbaoding2818/chainly/config"
)

var (
	ErrMsgQueueFull error = errors.New("msg queue full")
	flushInterval         = 1 * time.Minute // 最大等待间隔
)

type RabbitMQ struct {
	bufferChanel chan []byte
	batchSize    int32
	ticker       *time.Ticker
	// 通道关闭时 发送剩余的消息
	ctx context.Context
	wg  *sync.WaitGroup
}

func NewRabbitMQ(ctx context.Context, wg *sync.WaitGroup, cfg config.OperationLog, mgCfg config.MsgQueue) *RabbitMQ {
	mq := &RabbitMQ{
		ctx:          ctx,
		bufferChanel: make(chan []byte, cfg.BufferChanelSize),
		batchSize:    cfg.BatchSize,
		wg:           wg,
	}

	if cfg.BatchSize > 1 {
		// 1分钟超时 1分钟内如果没有累计到足够的消息，则强制发送
		mq.ticker = time.NewTicker(flushInterval)
		go mq.runBatch()
	} else {
		go mq.run()
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
				if r.batchSize > 1 { // 批量处理
					if len(msgs) < int(r.batchSize) {
						msgs = append(msgs, msg)
					} else {
						_msg := r.combineMessages(msgs)
						r.public(_msg)
						msgs = msgs[:0]
					}
				}
			case <-r.ticker.C:
				r.public(r.combineMessages(msgs))
			case <-r.ctx.Done():
				// 关闭通道时，将剩余的消息发送出去
				close(r.bufferChanel)
				r.public(r.combineMessages(msgs))
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
				r.public(msg)
			case <-r.ctx.Done():
				close(r.bufferChanel)
				// 关闭通道时，将剩余的消息发送出去
				for msg := range r.bufferChanel {
					r.public(msg)
				}
			}
		}
	}()
}

func (r *RabbitMQ) public(msgs []byte) {
	// 空数据检测
	if len(msgs) == 0 {
		return
	}

}

// combineMessages 合并多条消息为批量格式
func (r *RabbitMQ) combineMessages(msgs [][]byte) []byte {
	var buf bytes.Buffer
	for _, msg := range msgs {
		buf.Write(msg)
	}
	return buf.Bytes()
}
