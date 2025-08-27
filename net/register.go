package net

import (
	"sync"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
	"github.com/chenbaoding2818/chainly/net/actor"
	"github.com/chenbaoding2818/chainly/net/connection"
)

var (
	netComponent     *NetManager
	netComponentOnce sync.Once
)

type NetManager struct {
	cfg *config.Net
}

func (nm *NetManager) Start(cfg *config.Config) {
	nm.cfg = cfg.NetCfg
	// 检测websocket认证配置 为了安全考虑必须配置
	if nm.cfg.Ws.WsAuthHooker == nil {
		panic("websocket auth hooker is nil")
	}
	// 初始化actor模块
	actor.NewActorManager(nm.cfg)
	// 初始化连接管理器模块
	connection.NewConnManager()
}

func (nm *NetManager) Priority() uint32 {
	return lifecycle.NormalPriority
}

func (nm *NetManager) Stop() {

}

func NewSensitiveComponent() *NetManager {
	netComponentOnce.Do(func() {
		if netComponent == nil {
			netComponent = &NetManager{}
		}
	})

	return netComponent
}
