package iface

type MQType int8

const (
	// 默认使用rabbitmq
	MQTypeRabbitMQ MQType = iota + 1
	MQTypeKafka
)

type WaitForConfirmFunc func() error

// IProducer 生产者接口
type IProducer interface {
	// 同步发送消息 (等待确认)
	SendWithConfirm(msg []byte, f WaitForConfirmFunc) error
	// 异步发送消息 (不等待确认)
	SendWithoutConfirm(msg []byte) error
	// 合并消息 (批量发送时使用)
	CombineMessages(msgs [][]byte) []byte
	// 连接
	Connect() error
	// 重连
	Reconnect() error
}

// IConsumer 消费者接口
type IConsumer interface {
}
