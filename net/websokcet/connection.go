package websokcet

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chenbaoding2818/chainly/net/actor"
	"github.com/gorilla/websocket"
)

// 连接管理
type WsConnection struct {
	// // 连接ID
	// ID string
	// 玩家accountId
	AccountId string
	// 连接对象
	Conn *websocket.Conn
	// actor对象
	Actor  *actor.Actor
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

func (c *WsConnection) Listen() {
	go func() {
		for {
			select {
			// 关闭连接 例如正常退出、不同端顶号、服务异常等操作
			case <-c.ctx.Done():
				// TODO： 增加日志打印
				return
			default:
				// 设置心跳时间
				c.Conn.SetReadDeadline(time.Now().Add(time.Minute))
				// 读取信息
				msgType, msg, err := c.Conn.ReadMessage()
				if err != nil { // 读取异常触发关闭事件
					c.cancel()
				} else if msgType != websocket.BinaryMessage { // 数据类型不符，忽略本次信息
					continue
				} else { // 解析信息
					fmt.Printf("receive message: %s\n", string(msg))
					// 将消息发送给消息处理器
					c.Actor.SendMessage(msg)
				}
			}
		}
	}()
}

func (c *WsConnection) WriteMessage(msg []byte) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.Conn.WriteMessage(websocket.BinaryMessage, msg)
}
