package config

type LogMQ struct {
	// 消息队列的类型 支持的mq目前有1:rabbitmq 2:kafka
	MQType   int8
	RabbitMQ *RabbitMQConfig
	Kafka    *KafkaConfig
}

type BaseLog struct {
	// 是否开启远程
	RemoteEnabled bool
	// 缓冲区大小
	BufferChanelSize int32
	// 批量发送的大小 （开启会减少大量的网络io，但是会有丢失部分日志的风险，一般情况下是允许丢失一部分日志的）
	BatchSize int32
	LogMQ
}

// ServerLog 服务器日志配置
type ServerLog struct {
	// 日志的输出路径
	Path string
	// 是否要开启异步写磁盘方式
	AsyncEnabled bool
	BaseLog
}

// OperationLog 操作(运营)日志配置
type OperationLog struct {
	BaseLog
}

type Log struct {
	// 日志等级
	Level int
	// 写入磁盘的方式 0：同步 1：异步
	FWMode       int
	ServerLog    *ServerLog
	OperationLog *OperationLog
}
