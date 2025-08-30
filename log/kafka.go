package log

import (
	"context"
	"sync"

	"github.com/chenbaoding2818/chainly/config"
)

type Kafka struct {
	ctx context.Context
	wg  *sync.WaitGroup
}

func NewKafka(ctx context.Context, wg *sync.WaitGroup, cfg config.OperationLog, mgCfg config.MsgQueue) *Kafka {
	return &Kafka{
		ctx: ctx,
		wg:  wg,
	}
}

func (k *Kafka) SendMsg(msg []byte) error {
	// TODO: implement Kafka send message
	return nil
}
