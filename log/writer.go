package log

import (
	"io"

	"github.com/rs/zerolog"
)

type fileWriter struct {
}

func NewFileWriter() *fileWriter {
	return &fileWriter{}
}

func (w *fileWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type StdoutWriter struct {
	stdout io.Writer
}

func (w *StdoutWriter) Write(p []byte) (n int, err error) {
	return w.stdout.Write(p)
}

type RemoteWriter struct {
	batchSize int32
	channel   chan []byte
}

func (w *RemoteWriter) Write(p []byte) (n int, err error) {
	if w.batchSize > 0 {
		w.channel <- p
	}
	return len(p), nil
}

// BufferWriter 日志信息通过发送到channel进行缓存异步进行处理 优化日志写入性能
type BufferWriter struct {
	buffer      chan []byte
	multiWriter zerolog.LevelWriter
	batchSize   int32
}

func NewBufferWriter(bufferSize, batchSize int32, multiWriter zerolog.LevelWriter) *BufferWriter {
	return &BufferWriter{
		batchSize:   batchSize,
		multiWriter: multiWriter,
		buffer:      make(chan []byte, bufferSize),
	}
}

func (w *BufferWriter) run() {
	go func() {
		for p := range w.buffer {
			w.multiWriter.Write(p)
		}
	}()
}

func (w *BufferWriter) Write(p []byte) (n int, err error) {
	w.buffer <- p
	return len(p), nil
}
