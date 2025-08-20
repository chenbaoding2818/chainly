package websokcet

import (
	"context"

	"github.com/gorilla/websocket"
)

// 连接管理
type WsConnection struct {
	// 连接ID
	ID string
	// 连接对象
	Conn *websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc
}
