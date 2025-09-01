package mq

import (
	"context"
	"sync"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
)

var (
	mqComponent     *MQManager
	mqComponentOnce sync.Once
)

type MQManager struct {
}

func (lm *MQManager) Start(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) {

}

func (lm *MQManager) Priority() int32 {
	return lifecycle.HighPriority
}

func (lm *MQManager) Stop() {

}

func NewMQComponent() *MQManager {
	mqComponentOnce.Do(func() {
		if mqComponent == nil {
			mqComponent = &MQManager{}
		}
	})
	return mqComponent
}
