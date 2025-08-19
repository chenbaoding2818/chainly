package log

import "sync"

var (
	LogComponent     *LogManager
	LogComponentOnce sync.Once
)

type LogManager struct {
}

func NewLogComponent() *LogManager {
	LogComponentOnce.Do(func() {
		if LogComponent == nil {
			LogComponent = &LogManager{}
		}
	})

	return LogComponent
}
