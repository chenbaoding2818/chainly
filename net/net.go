package net

import (
	"sync"

	"github.com/chenbaoding2818/chainly/config"
)

var (
	netComponent     *NetManager
	netComponentOnce sync.Once
)

type NetManager struct {
	cfg *config.Net
}

func NewSensitiveComponent() *NetManager {
	netComponentOnce.Do(func() {
		if netComponent == nil {
			netComponent = &NetManager{}
		}
	})

	return netComponent
}

// Run 运行网络模块
func (nm *NetManager) Run() {

}
