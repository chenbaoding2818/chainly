package log

import "github.com/chenbaoding2818/chainly/config"

type RabbitMQ struct {
	BufferChanel chan []byte
	BatchSize    int32
}

func NewRabbitMQ(cfg config.OperationLog, mgCfg config.MsgQueue) *RabbitMQ {
	return &RabbitMQ{
		BufferChanel: make(chan []byte, cfg.BufferChanelSize),
		BatchSize:    cfg.BatchSize,
	}
}

func (r *RabbitMQ) SendMsg(msg []byte) error {
	return nil
}
