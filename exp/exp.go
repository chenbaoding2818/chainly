package exp

import (
	"fmt"

	iface "github.com/chenbaoding2818/chainly/interface"
	"github.com/chenbaoding2818/chainly/net/actor"
	"github.com/chenbaoding2818/chainly/net/connection"
)

// 怎么使用？
type GameWsHandler struct {
	handlermap map[int]func(actor *actor.Actor, msg []byte) error
}

func (gh *GameWsHandler) Handle(connPtr uintptr, msg []byte) error {
	// 反序列化msg变成游戏定义的消息
	// 伪代码：
	// 1、根据msg的cmdid,找到对应的处理函数
	// 2、调用处理函数处理消息
	// 3、返回处理结果 如果有error，则返回error 调用底层返回错误信息

	// 这样的好处是：1、可以统一记录每个协议的操作记录
	actor := connection.NewConnManager().GetActor(connPtr)

	fmt.Println(actor)

	return nil
}

func NewGamweHandler() iface.IMessageHandler {
	return &GameWsHandler{}
}
