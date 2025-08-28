package iface

type LogLevel int8

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled
)

type ILog interface {
	SetLevel(level LogLevel)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Fatal(msg string)
	Panic(msg string)
	Error(msg string)
	// Report专门上报行为日志、运营日志
	Report(msg []byte)
}

type ILogProducer interface {
	SendMsg(msg []byte) error
	// SendMsgBatch(msgs [][]byte) error
}
