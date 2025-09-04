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
	// dst目标名称，例如rabbitmq的routingkey名称，kafka的topic名称
	SendWithConfirm(dst string, msg []byte, successCallback, failCallback WaitForConfirmFunc) error
	// 异步发送消息 (不等待确认)
	SendWithoutConfirm(dst string, msg []byte) error
	// 合并消息 (批量发送时使用)
	CombineMessages(msgs [][]byte) []byte
	// 连接
	Connect() error
	// 重连
	Reconnect() error
}

// ConsumerHandlerFunc 消费者处理函数类型
type ConsumerHandlerFunc func([]byte) error

// IConsumer 消费者接口
type IConsumer interface {
}
