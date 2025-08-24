package actor

import (
	"fmt"
	"runtime/debug"
	"time"

	iface "github.com/chenbaoding2818/chainly/interface"
	"github.com/gorilla/websocket"
)

var idleDuration = 10 * time.Minute

// Actor 每个玩家持有一个Actor对象，用于处理消息
// 什么时候回收Actor对象？
type Actor struct {
	// 玩家ID
	Id string
	// 连接
	conn *websocket.Conn
	// handler 处理消息的函数
	handler iface.IMessageHandler
	// 消息队列(信箱)
	mailbox chan []byte
	// 退出信号
	quitCh chan struct{}
	// 离线定时器 离线定时器（用于在离线一段时间后销毁自己）
	// offlineTimer
	//
	isOnline bool
}

func (a *Actor) Start() {
	defer func() {
		// TODO: 处理panic 因为在处理消息时，可能出现panic，需要处理
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	// 设置一个定时器，用于在离线一段时间后自动退出
	idleTimer := time.NewTimer(idleDuration)
	defer idleTimer.Stop()

	for {
		select {
		case msg := <-a.mailbox:
			// TODO：处理消息
			fmt.Printf("Actor %s received message: %v\n", a.Id, msg)
			// 如果处理笑消息有问题，则需要处理错误
			a.processMessage(msg)

			// if !idleTimer.Stop() {
			// 	<-idleTimer.C
			// }
			// idleTimer.Reset(idleDuration)
		case <-idleTimer.C:
			// // 一段时间没有消息（且玩家离线），则退出Actor
			// if !p.online {
			// 	// 从管理器中移除
			// 	playerManager.actors.Delete(p.ID)
			// 	return
			// }
		case <-a.quitCh:
			// 退出
			return
		}
	}
}

func (a *Actor) SendMessage(msg []byte) {
	a.mailbox <- msg
}

// processMessage 处理信箱中的消息
func (a *Actor) processMessage(msg []byte) error {
	return a.handler.Handle(msg)
}

func (a *Actor) Stop() {
	// 关闭信箱
	close(a.mailbox)
	// 关闭退出信号
	close(a.quitCh)
	// 同时通知管理器删除该Actor
}

func (a *Actor) GetConn() *websocket.Conn {
	return a.conn
}

func (a *Actor) IsOnline() bool {
	return a.isOnline
}
