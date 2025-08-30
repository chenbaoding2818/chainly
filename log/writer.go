package log

import (
	"context"
	"io"
	"sync"

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

func NewStdoutWriter(stdout io.Writer) *StdoutWriter {
	return &StdoutWriter{
		stdout: stdout,
	}
}

func (w *StdoutWriter) Write(p []byte) (n int, err error) {
	return w.stdout.Write(p)
}

// BufferWriter 日志信息通过发送到channel进行缓存异步进行处理 优化日志写入性能
type BufferWriter struct {
	buffer      chan []byte
	multiWriter zerolog.LevelWriter
	batchSize   int32
	ctx         context.Context
	wg          *sync.WaitGroup
}

func NewBufferWriter(ctx context.Context, wg *sync.WaitGroup, bufferSize, batchSize int32,
	multiWriter zerolog.LevelWriter) *BufferWriter {
	return &BufferWriter{
		batchSize:   batchSize,
		multiWriter: multiWriter,
		ctx:         ctx,
		wg:          wg,
		buffer:      make(chan []byte, bufferSize),
	}
}

func (w *BufferWriter) run() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				// 消费完所有日志信息后退出
				close(w.buffer)
				for p := range w.buffer {
					w.multiWriter.Write(p)
				}
				return
			case p := <-w.buffer:
				w.multiWriter.Write(p)
			}
		}
	}()

}

func (w *BufferWriter) Write(p []byte) (n int, err error) {
	w.buffer <- p
	return len(p), nil
}
