package websokcet

import (
	"github.com/gin-gonic/gin"
)

func websocketUpgrader(ctx *gin.Context) {
	// WRANNING: 严格来说应该在websocket握手前进行鉴权以提供更好的性能与安全性
	// 升级服务
	// conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	// if err != nil {
	// 	log.Error(fmt.Sprintf("websocket upgrade error, error:%s", err.Error()))
	// 	return
	// }
	// // 验证入参
	// vaild := &Vaild{}
	// if err := ctx.ShouldBind(vaild); err != nil {
	// 	log.Warn(fmt.Sprintf("params not vaild, error:%s", err.Error()))
	// 	PushConnetError(conn, pb.ErrorCode_EC_VaildFail, "")
	// 	return
	// } else if config.App.CloseServer {
	// 	PushConnetError(conn, pb.ErrorCode_EC_ServerMaintenance)
	// 	return
	// }
	// // 解析Token
	// resp, err := grpc_client.BackendSyncRequest[*pb.GetRoleByTokenResp](pb.SyncCmd_SC_GetRoleByToken, &pb.GetRoleByTokenReq{
	// 	Token:      vaild.Token,
	// 	PlatformId: config.App.PlatformId,
	// 	ServerId:   config.App.ServerId,
	// })
	// if err != nil {
	// 	log.Error(fmt.Sprintf("request backend server error,Token:%s  PlatformId:%s ServerId:%d  vaild:%v error:%s", vaild.Token, config.App.PlatformId, config.App.ServerId, vaild, err.Error()))
	// 	PushConnetError(conn, pb.ErrorCode_EC_ServerError)
	// 	return
	// }
	// // 判断是区服状态和否黑名单
	// if resp.Status == pb.PlayerStatus_BlackList && general.NowMilli() < resp.BlackListEnd { // 黑名单
	// 	PushConnetError(conn, pb.ErrorCode_EC_AccountDisable, resp.Desc, general.Now().Format(time.DateTime))
	// 	return
	// } else if resp.ServerStatus == pb.PlayerServerStatus_Service_Shutdown_Maintenance &&
	// 	resp.Status != pb.PlayerStatus_WhiteList { // 服务维护中
	// 	PushConnetError(conn, pb.ErrorCode_EC_ServerMaintenance)
	// 	return
	// }

	// account := fmt.Sprintf("%s_%d", resp.Account, vaild.ServerId)
	// playerIdStr, err := cache.DB.HGet(cache.Accounts, account)
	// if err != nil && err != cache.ErrFooNil {
	// 	PushConnetError(conn, pb.ErrorCode_EC_ServerNetworkError)
	// 	log.Error(fmt.Sprintf("[client][WebsocketUpgrader]account(%s) fail to get player_id by account, error:%s", resp.Account, err.Error()))
	// 	return
	// }
	// playerId := int32(0)
	// if playerIdStr != "" {
	// 	playerId, err = utils.StrToNumber[int32](playerIdStr)
	// 	if err != nil {
	// 		PushConnetError(conn, pb.ErrorCode_EC_ServerNetworkError)
	// 		log.Error(fmt.Sprintf("[client][WebsocketUpgrader]playerId(%s) fail to change to number, error:%s", playerIdStr, err.Error()))
	// 		return
	// 	}
	// }
	// // 玩家会话
	// playerSession := &Session{
	// 	Conn:      conn,
	// 	Reconnet:  make(chan bool),
	// 	Disconnet: make(chan bool),
	// 	Msg:       make(chan *pb.Msg),
	// 	Account:   account,
	// 	PlayerId:  playerId,
	// 	Lock:      &sync.Mutex{},
	// }
	// playerSession.RegisterSession()
	// playerSession.Listen()
}
