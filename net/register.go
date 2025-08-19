package net

import (
	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
)

func (nm *NetManager) Start(cfg *config.Config) {
	nm.cfg = cfg.NetCfg
	nm.Run()
}

func (nm *NetManager) Priority() uint32 {
	return lifecycle.NormalPriority
}

func (nm *NetManager) Stop() {

}
