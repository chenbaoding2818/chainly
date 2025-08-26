package websokcet

import (
	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
	"github.com/chenbaoding2818/chainly/net/actor"
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

func (ws *WebsocketService) websocketUpgrader(ctx *gin.Context) {
	var (
		wsCfg    = ws.cfg.Ws
		wsConn   *websocket.Conn
		authResp iface.IWsAuthResp
		conn     iface.IConnection
		err      error
	)
	// 检测服务最大连接数
	if connection.NewConnManager().Len() > ws.cfg.MaxConn {
		// TODO: 增加日志打印Error
		return
	}
	// 进行websocket认证 必须进行认证
	if wsCfg.WsAuthHooker == nil {
		// TODO: 增加日志打印Error
		return
	}
	authResp, err = wsCfg.WsAuthHooker(ctx.Request)
	if err != nil {
		// TODO: 增加日志打印Error
		ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}
	// 协议升级
	wsConn, err = wsCfg.Upgrader.Upgrade(ctx.Writer, ctx.Request, wsCfg.WsRespHeader)
	if err != nil {
		// TODO: 增加日志打印Error
		return
	}
	// 处理每个连接
	conn = NewWsConnection(wsConn)
	if authResp != nil {
		// 设置玩家账号id信息
		conn.SetAccountId(authResp.GetAccountId())
		// 设置玩家id信息
		conn.SetplayerId(authResp.GetPlayerId())
	}
	// 检测连接中玩家的信息

	// 创建一个actor 每个连接对应一个actor
	actor := actor.NewActor(conn, ws.cfg)
	// 连接管理器添加新的连接
	connection.NewConnManager().AddConn(actor)
	// actor开始监听连接信息
	actor.ListenConn()
}
