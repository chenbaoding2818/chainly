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
	// 创建生产者
	NewProducer(ctx, wg, cfg)
	// 创建消费者 TODO: 带实现

}

func (lm *MQManager) Priority() int32 {
	return lifecycle.NormalPriority
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
