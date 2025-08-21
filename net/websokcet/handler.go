package websokcet

import (
	"fmt"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/gin-gonic/gin"
)

type WebsocketService struct {
	cfg *config.Net
}

func (ws *WebsocketService) Route(route *gin.Engine) {
	route.GET(ws.cfg.Ws.WsPath, ws.websocketUpgrader)
}

func (ws *WebsocketService) websocketUpgrader(ctx *gin.Context) {
	// 检测服务最大连接数 TODO: 如何判断？连接池管理进行判断？
	if 100 > ws.cfg.MaxConn {
		// TODO: 增加日志打印Error
		return
	}

	if ws.cfg.Ws.WsAuthHooker != nil {
		err := ws.cfg.Ws.WsAuthHooker(ctx.Request)
		if err != nil {
			// TODO: 增加日志打印Error
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
	}
	// 协议升级
	conn, err := ws.cfg.Ws.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// TODO: 增加日志打印Error
		return
	}
	fmt.Printf("websocket connection established: %s\n", conn.RemoteAddr().String())
	wsConn := NewWsConnection("", conn)
	go wsConn.StartReader()
	// // 加入连接管理器
	// actor.NewActorManager().GetActor("websocket")
}
