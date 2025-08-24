package iface

// IMesssageHandler 处理消息的接口
// 用户只要实现这个接口，就可以自定义自己的消息处理逻辑
type IMessageHandler interface {
	Handle(msg []byte) error
}
