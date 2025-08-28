package log

import (
	"github.com/rs/zerolog"
)

type callerHook struct {
	callerSkipFrameCount int32
}

func (ch callerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.CallerSkipFrame(int(ch.callerSkipFrameCount))
	e.Caller(int(ch.callerSkipFrameCount))
}

func NewDefaultCallerHook() zerolog.Hook {
	return callerHook{callerSkipFrameCount: 2}
}

type RemoterHook struct {
}

func (rh RemoterHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.NoLevel { // 运营业务的日志 即时上传
		// if err := rh.producer.SendMsg([]byte(msg)); err != nil {
		// 	fmt.Println("rabbitmq hook send msg error:", err)
		// }
		e.Discard()
	}
}
