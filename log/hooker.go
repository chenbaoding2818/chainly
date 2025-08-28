package log

import (
	"github.com/chenbaoding2818/chainly/config"
	"github.com/rs/zerolog"
)

type callerHook struct {
	callerSkipFrameCount int32
}

func NewDefaultCallerHook() zerolog.Hook {
	return callerHook{callerSkipFrameCount: 2}
}

func (ch callerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.CallerSkipFrame(int(ch.callerSkipFrameCount))
	e.Caller(int(ch.callerSkipFrameCount))
}

type RemoterHook struct {
}

func NewRemoterHook(cfg config.OperationLog) zerolog.Hook {
	return RemoterHook{}
}

func (rh RemoterHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// 运营业务的日志，发送到消息队列
	if level == zerolog.NoLevel {
		// if err := rh.producer.SendMsg([]byte(msg)); err != nil {
		// 	fmt.Println("rabbitmq hook send msg error:", err)
		// }
		e.Discard()
	}
}
