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
					break
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
