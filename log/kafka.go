package log

import (
	"context"

	"github.com/chenbaoding2818/chainly/config"
)

type Kafka struct {
}

func NewKafka(ctx context.Context, cfg config.OperationLog, mgCfg config.MsgQueue) *Kafka {
	return &Kafka{}
}

func (k *Kafka) SendMsg(msg []byte) error {
	return nil
}
