package log

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"

	"github.com/rs/zerolog"
)

var (
	logger     iface.ILog
	loggerOnce sync.Once
)

type Logger struct {
	name     string
	serverId int
	level    zerolog.Level
	log      zerolog.Logger
}

func NewLogger(config *config.Config) iface.ILog {
	loggerOnce.Do(func() {
		if logger == nil {
			writers := make([]io.Writer, 0)
			if len(config.LogCfg.ServerLog.Path) > 0 {
				writers = append(writers, NewFileWriter())
			} else {
				writers = append(writers, &StdoutWriter{os.Stdout})
			}
			// 是否开启远程日志
			if config.LogCfg.ServerLog.RemoteEnabled {
				writers = append(writers, &RemoteWriter{})
			}
			var multiWriter io.Writer
			// 系统日志是否开启异步写入
			if config.LogCfg.ServerLog.AsyncEnabled {
				multiW := zerolog.MultiLevelWriter(writers...)
				bufferWriter := NewBufferWriter(config.LogCfg.ServerLog.BufferChanelSize,
					config.LogCfg.ServerLog.BatchSize,
					multiW)
				multiWriter = bufferWriter.multiWriter
				go bufferWriter.run()
			} else {
				multiWriter = zerolog.MultiLevelWriter(writers...)
			}
			// 设置全局日志等级
			zerolog.SetGlobalLevel(zerolog.Level(config.LogCfg.Level))
			// 设置日志时间格式
			zerolog.TimeFieldFormat = time.RFC3339
			l := zerolog.New(multiWriter)
			logger = &Logger{
				name:     config.BasicCfg.GameName,
				serverId: config.BasicCfg.ServerId,
				log:      l.With().Timestamp().Logger().Hook(nil),
			}
		}
	})
	return logger
}

// SetLevel 设置日志等级 预留api以供动态调整日志等级
func (l *Logger) SetLevel(level iface.LogLevel) {
	zerolog.SetGlobalLevel(zerolog.Level(level))
}

func (l *Logger) Debug(msg string) {
	l.log.Debug().Str("gameName", l.name).Int("serverId", l.serverId).Msg(msg)

}

func (l *Logger) Info(msg string) {
	l.log.Info().Str("gameName", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Warn(msg string) {
	l.log.Warn().Str("gameName", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Error(msg string) {
	l.log.Error().Str("gameName", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Fatal(msg string) {
	l.log.Fatal().Str("gameName", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Panic(msg string) {
	l.log.Panic().Str("gameName", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Report(msg []byte) {
	l.log.Log().Msg(string(msg))
}
