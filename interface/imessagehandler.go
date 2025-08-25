package iface

// IMesssageHandler 处理消息的接口
// 用户只要实现这个接口，就可以自定义自己的消息处理逻辑
type IMessageHandler interface {
	Handle(msg []byte) error
}

// 怎么使用？
type GameHandler struct{}

func (gh *GameHandler) Handle(msg []byte) error {
	// 反序列化msg变成游戏定义的消息
	// 伪代码：
	// 1、根据msg的cmdid,找到对应的处理函数
	// 2、调用处理函数处理消息
	// 3、返回处理结果 如果有error，则返回error 调用底层返回错误信息

	return nil
}

func NewGamweHandler() IMessageHandler {
	return &GameHandler{}
}
