package net

import (
	"sync"

	"github.com/chenbaoding2818/chainly/config"
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
	// 检测websocket认证配置 为了安全考虑必须配置
	if nm.cfg.Ws.WsAuthHooker == nil {
		panic("websocket auth hooker is nil")
	}
	// 初始化actor模块
	actor.NewActorManager(nm.cfg)
	// 初始化连接管理器模块
	connection.NewConnManager()
}
