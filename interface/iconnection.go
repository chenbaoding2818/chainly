package iface

import (
	"context"
	"time"
)

type IConnection interface {
	// // 监听连接的消息
	// Listen()
	// 设置连接的玩家账号ID
	SetAccountId(accountId string)
	SetplayerId(playerId string)
	// 获取连接的玩家账号ID
	GetAccountId() string
	GetConnPtr() uintptr
	ReadMessage() (int, []byte, error)
	// 向连接发送消息
	WriteMessage(msg []byte) error
	//
	GetCtx() context.Context
	GetCancel() context.CancelFunc
	SetReadDeadline(t time.Time) error
	Close() error
}
