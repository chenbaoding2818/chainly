package websokcet

import (
	"unsafe"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
	"github.com/chenbaoding2818/chainly/net/connection"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	cfg *config.Net
}

func (ws *WebsocketService) Route(route *gin.Engine) {
	route.GET(ws.cfg.Ws.WsPath, ws.websocketUpgrader)
}

// account id关联 每个连接WSconn key:accountId value:WsConnection
// 每个连接地址关联一个每个连接WSconn WS指向一个actor key:connAddr value:WsConnection
func (ws *WebsocketService) websocketUpgrader(ctx *gin.Context) {
	var (
		wsCfg    = ws.cfg.Ws
		conn     *websocket.Conn
		authResp iface.IWsAuthResp
		iconn    iface.IConnection
		err      error
	)
	// 检测服务最大连接数 TODO: 如何判断？连接池管理进行判断？
	if connection.NewConnManager().Len() > ws.cfg.MaxConn {
		// TODO: 增加日志打印Error
		return
	}
	if wsCfg.WsAuthHooker != nil {
		authResp, err = wsCfg.WsAuthHooker(ctx.Request)
		if err != nil {
			// TODO: 增加日志打印Error
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
	}
	// 协议升级
	conn, err = wsCfg.Upgrader.Upgrade(ctx.Writer, ctx.Request, wsCfg.WsRespHeader)
	if err != nil {
		// TODO: 增加日志打印Error
		return
	}
	// 获取改连接的地址
	ptr := uintptr(unsafe.Pointer(conn))

	iconn = NewWsConnection(conn)
	connection.NewConnManager().AddConn(ptr, iconn)
	if authResp != nil {
		authResp.GetAccount()
	}

	iconn.Listen()
	// // 加入连接管理器
	// actor.NewActorManager().GetActor("websocket")
}
