package iface

import "context"

type MQType int8

const (
	// 默认使用rabbitmq
	MQTypeRabbitMQ MQType = iota + 1
	MQTypeKafka
)

// IProducer 生产者接口
type IProducer interface {
	// 同步发送消息
	Send(ctx context.Context, msg []byte) error
	// 异步发送消息
	SendAsync(ctx context.Context, msg []byte) error
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
