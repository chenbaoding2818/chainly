package actor

import (
	"fmt"
	"runtime/debug"
	"time"

	iface "github.com/chenbaoding2818/chainly/interface"
	"github.com/gorilla/websocket"
)

var idleDuration = 10 * time.Minute

// Actor 每个连接/玩家拥有一个Actor对象，用于处理消息 保证数据并发安全
// 什么时候回收Actor对象？
type Actor struct {
	// 玩家ID
	Id string
	// 连接
	Conn iface.IConnection
	// handler 处理消息的函数
	handler iface.IMessageHandler
	// 消息队列(信箱) 玩家所有需要处理的消息全部进入信箱才能被处理 保证数据并发安全
	mailbox chan []byte
	// 退出信号
	quitCh chan struct{}
	// 离线定时器 离线定时器（用于在离线一段时间后销毁自己）
	// offlineTimer
	//
	isOnline bool
	// 断开连接需要处理消息
	disconnecttionMsg []byte
	// 异端登陆时，旧连接需要断开

}

// SetDisconnectMsg 设置连接断开时需要处理的消息
func (a *Actor) SetDisconnectMsg(msg []byte) {
	a.disconnecttionMsg = msg
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

func (a *Actor) ConnClose() {
	a.Conn.Close()
	a.Conn = nil
}

func (a *Actor) SendMessage(msgType int, msg []byte, err error) {
	// 读取异常触发关闭事件
	if err != nil {
		// 连接关闭
		a.ConnClose()
		// 连接关闭就是玩家登出
		// 将玩家登出的消息放入信箱进行串行处理，保证数据并发安全
		msg = a.disconnecttionMsg
	} else if msgType != websocket.BinaryMessage {
		return
	}

	a.mailbox <- msg
}

// processMessage 处理信箱中的消息 要传入玩家信息 怎么传入玩家信息
func (a *Actor) processMessage(msg []byte) error {
	err := a.handler.Handle(a.Conn.GetConnPtr(), msg)
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

func (a *Actor) ListenConn() {
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
				a.SendMessage(a.Conn.ReadMessage())
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
	// 开启信箱监听
	go actor.start()
	return actor
}
