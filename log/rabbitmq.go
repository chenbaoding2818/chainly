package log

import (
	"bytes"
	"context"
	"errors"
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
	timer        *time.Timer
	// 通道关闭时 发送剩余的消息
	ctx context.Context
}

func NewRabbitMQ(ctx context.Context, cfg config.OperationLog, mgCfg config.MsgQueue) *RabbitMQ {
	mq := &RabbitMQ{
		ctx:          ctx,
		bufferChanel: make(chan []byte, cfg.BufferChanelSize),
		batchSize:    cfg.BatchSize,
	}

	if cfg.BatchSize > 1 {
		// 1分钟超时 1分钟内如果没有累计到足够的消息，则强制发送
		mq.timer = time.NewTimer(flushInterval)
		go mq.runBatch()
	} else {

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
	go func() {
		msgs := make([][]byte, 0, r.batchSize)
		for {
			select {
			case msg := <-r.bufferChanel:
				if r.batchSize > 1 { // 批量处理
					if len(msgs) < int(r.batchSize) {
						msgs = append(msgs, msg)
					} else {
						msgs = append(msgs, msg)
						r.public(msgs)
						msgs = msgs[:0]
					}
				}
			case <-r.timer.C:
				r.public(msgs)
			case <-r.ctx.Done():
				r.public(msgs)
				return
			}
		}
	}()
}

// run 运行消费协程
func (r *RabbitMQ) run() {
	go func() {
		for {
			select {
			case msg := <-r.bufferChanel:
				r.public([][]byte{msg})
			}
		}
	}()
}

func (r *RabbitMQ) public(msgs [][]byte) {

}

// combineMessages 合并多条消息为批量格式
func (r *RabbitMQ) combineMessages(msgs [][]byte) []byte {
	var buf bytes.Buffer
	for _, msg := range msgs {
		buf.Write(msg)
	}
	return buf.Bytes()
}
