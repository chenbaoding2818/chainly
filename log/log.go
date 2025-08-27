package log

import (
	"os"
	"sync"

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
	log      zerolog.Logger
}

func NewLogger(config *config.Config) iface.ILog {
	loggerOnce.Do(func() {
		if logger == nil {
			logger = &Logger{
				name:     config.BasicCfg.GameName,
				serverId: config.BasicCfg.ServerId,
				log: zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
					With().Timestamp().Logger().Hook(nil),
			}
		}
	})
	return logger
}

func (l *Logger) SetLevel(level iface.LogLevel) {
	zerolog.SetGlobalLevel(zerolog.Level(level))
}

func (l *Logger) Debug(msg string) {
	l.log.Debug().Str("name", l.name).Int("serverId", l.serverId).Msg(msg)

}

func (l *Logger) Info(msg string) {
	l.log.Info().Str("name", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Warn(msg string) {
	l.log.Warn().Str("name", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Error(msg string) {
	l.log.Error().Str("name", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Fatal(msg string) {
	l.log.Fatal().Str("name", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Panic(msg string) {
	l.log.Panic().Str("name", l.name).Int("serverId", l.serverId).Msg(msg)
}

func (l *Logger) Report(msg []byte) {
	l.log.Log()
}
