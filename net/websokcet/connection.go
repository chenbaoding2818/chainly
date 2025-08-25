package websokcet

import (
	"context"
	"sync"
	"time"
	"unsafe"

	"github.com/gorilla/websocket"
)

// 连接管理
type WsConnection struct {
	// // 连接ID
	// ID string
	// 玩家accountId
	AccountId string
	// 连接对象
	Conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
	// 写锁 用于写消息
	Lock *sync.Mutex
}

func NewWsConnection(conn *websocket.Conn) *WsConnection {
	ctx, cancel := context.WithCancel(context.Background())
	return &WsConnection{
		// ID:     id,
		Conn:   conn,
		ctx:    ctx,
		cancel: cancel,
		Lock:   new(sync.Mutex),
	}
}

func (c *WsConnection) GetCancel() context.CancelFunc {
	return c.cancel
}

func (c *WsConnection) GetCtx() context.Context {
	return c.ctx
}

func (c *WsConnection) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *WsConnection) ReadMessage() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (c *WsConnection) SetAccountId(accountId string) {
	c.AccountId = accountId
}

func (c *WsConnection) GetAccountId() string {
	return c.AccountId
}

func (c *WsConnection) GetConnPtr() uintptr {
	return uintptr(unsafe.Pointer(c.Conn))
}

func (c *WsConnection) Close() error {
	return c.Conn.Close()
}

func (c *WsConnection) WriteMessage(msg []byte) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.Conn.WriteMessage(websocket.BinaryMessage, msg)
}
