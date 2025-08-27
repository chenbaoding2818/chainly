package log

import "fmt"

func Debug(msg string) {
	// implementation of Debugf
	logger.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	// implementation of Debugf
	msg := fmt.Sprintf(format, args...)
	logger.Debug(msg)
}

func Info(msg string) {
	// implementation of Debugf
	logger.Info(msg)
}

func Infof(format string, args ...interface{}) {
	// implementation of Debugf
	msg := fmt.Sprintf(format, args...)
	logger.Info(msg)
}

func Warn(msg string) {
	// implementation of Debugf
	logger.Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	// implementation of Debugf
	msg := fmt.Sprintf(format, args...)
	logger.Warn(msg)
}

func Error(msg string) {
	// implementation of Debugf
	logger.Error(msg)
}

func Errorf(format string, args ...interface{}) {
	// implementation of Debugf
	msg := fmt.Sprintf(format, args...)
	logger.Error(msg)
}

func Fatal(msg string) {
	// implementation of Debugf
	logger.Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	// implementation of Debugf
	msg := fmt.Sprintf(format, args...)
	logger.Fatal(msg)
}

func Panic(msg string) {
	// implementation of Debugf
	logger.Panic(msg)
}

func Panicf(format string, args ...interface{}) {
	// implementation of Debugf
	msg := fmt.Sprintf(format, args...)
	logger.Panic(msg)
}

// Report 专门用于上报行为日志以及运营日志
func Report(msg []byte) {
	// implementation of Debugf
	logger.Report(msg)
}
