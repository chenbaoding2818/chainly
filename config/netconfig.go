package config

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Net struct {
	// 高并发自动扩容
	// 开启后，当为玩家在线人数处于高峰值时，实现channel自动扩容
	ConcurrencyAutoEnabe bool
	// 处理器超时开关
	TimeoutEmabled bool
	// 最大连接数
	MaxConn int

	Ws   *Websocket
	Kcp  *Kcp
	Http *Http
}

type Websocket struct {
	WsPort int
	// websocket连接升级器
	Upgrader *websocket.Upgrader
	// websocket连接认证hooker
	// 可以通过创建服务时传入options参数设置
	WsAuthHooker func(req *http.Request) error
}

func NewDefaultWebsocketCfg() *Websocket {
	return &Websocket{
		WsPort: 8000,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type Kcp struct {
	// kcp连接端口
	KcpPort int
}

func NewDefaultKcpCfg() *Kcp {
	return &Kcp{
		KcpPort: 8001,
	}
}

type Http struct {
	// http服务端口
	HttpPort int
}

// NewNetCfg 创建一个默认的Net配置
func NewDefaultNetCfg() *Net {
	return &Net{
		ConcurrencyAutoEnabe: false,
		TimeoutEmabled:       false,
		// 默认最大连接数
		MaxConn: 5000,
		Ws: &Websocket{
			Upgrader: &websocket.Upgrader{
				ReadBufferSize:  4096,
				WriteBufferSize: 4096,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	}
}
