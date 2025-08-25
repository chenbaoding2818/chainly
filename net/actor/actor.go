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
	Conn iface.IConnection
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

func (a *Actor) start() {
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

// processMessage 处理信箱中的消息 要传入玩家信息 怎么传入玩家信息
func (a *Actor) processMessage(msg []byte) error {
	err := a.handler.Handle(msg)
	if err != nil {
		a.Conn.WriteMessage([]byte(err.Error()))
	}
	return err
}

func (a *Actor) Stop() {
	// 关闭信箱
	close(a.mailbox)
	// 关闭退出信号
	close(a.quitCh)
	// 同时通知管理器删除该Actor
}

func (a *Actor) GetConn() iface.IConnection {
	return a.Conn
}

func (a *Actor) IsOnline() bool {
	return a.isOnline
}

func (a *Actor) Listen() {
	go func() {
		for {
			select {
			// 关闭连接 例如正常退出、不同端顶号、服务异常等操作
			case <-a.Conn.GetCtx().Done():
				// TODO： 增加日志打印
				return
			default:
				// 设置心跳时间
				a.Conn.SetReadDeadline(time.Now().Add(time.Minute))
				// 读取信息
				msgType, msg, err := a.Conn.ReadMessage()
				if err != nil { // 读取异常触发关闭事件
					break
				} else if msgType != websocket.BinaryMessage { // 数据类型不符，忽略本次信息
					continue
				} else { // 解析信息
					fmt.Printf("receive message: %s\n", string(msg))
					// 将消息发送给消息处理器
					a.SendMessage(msg)
				}
			}
		}
	}()
}

func NewActor(conn iface.IConnection) *Actor {
	actor := &Actor{
		Conn:    conn,
		mailbox: make(chan []byte, 100),
		quitCh:  make(chan struct{}),
	}
	go actor.start()
	return actor
}
