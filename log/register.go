package log

import (
	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
)

func (lm *LogManager) Start(cfg *config.Config) {

}

func (lm *LogManager) Priority() int32 {
	return lifecycle.HighPriority
}

func (lm *LogManager) Stop() {

}
