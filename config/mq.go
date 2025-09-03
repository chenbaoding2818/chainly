package config

// MsgQueue 定义消息队列配置，消息队列在游戏业务中应用广泛，以下是一些主要应用场景
// 1、系统解耦
// 2、异步操作处理
// 3、日志处理
// 4、流量削峰

type RabbitMQConfig struct {
	Urls []string
	// 交换机名称
	ExchangeName string
	// 交换机类型
	ExchangeType string
	// 路由键
	RoutingKey string
	// confirm模式是否开启
	ConfirmEnable bool
	// 确认超时时间
	MsgConfirmTimeout int
	// 消息持久化
	Persistent bool
	// 重连次数
	ReconnectCount int
	// 最大重连次数
	ReconnectMax int
	// mq心跳检测间隔
	HeartBeat int
	// 最大未确认消息数
	MaxPendingMessage int
}

type KafkaConfig struct {
}

type MsgQueue struct {
	RabbitMQ *RabbitMQConfig
	Kafka    *KafkaConfig
}
