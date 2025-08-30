package log

import (
	"context"
	"sync"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
)

var (
	LogComponent     *LogManager
	LogComponentOnce sync.Once
)

type LogManager struct {
}

func (lm *LogManager) Start(ctx context.Context, cfg *config.Config) {
	NewLogger(ctx, cfg)
}

func (lm *LogManager) Priority() int32 {
	return lifecycle.HighPriority
}

func (lm *LogManager) Stop() {

}

func NewLogComponent() *LogManager {
	LogComponentOnce.Do(func() {
		if LogComponent == nil {
			LogComponent = &LogManager{}
		}
	})
	return LogComponent
}
