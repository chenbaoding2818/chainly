package core

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
	"github.com/chenbaoding2818/chainly/log"
	"github.com/chenbaoding2818/chainly/sensitive"
	"github.com/chenbaoding2818/chainly/timer/cron"
)

type Server struct {
	cfg *config.Config
}

func NewDefaultServer(opts ...Option) *Server {
	return NewServer(config.NewDefaultConfig(), opts...)
}

func NewServer(cfg *config.Config, opts ...Option) *Server {
	for _, opt := range opts {
		opt(cfg)
	}
	return &Server{
		cfg: cfg,
	}
}

// RegisterDefaultComponents 注册默认服务组件
func (s *Server) RegisterDefaultComponents() {
	// 注册日志组件
	lifecycle.LifecycleMgr.AddLifecycle(log.NewLogComponent())
	// 注册敏感词库组件
	lifecycle.LifecycleMgr.AddLifecycle(sensitive.NewSensitiveComponent())
	// 注册定时组件
	lifecycle.LifecycleMgr.AddLifecycle(cron.NewCronComponent())
}

// RegisterComponent 注册服务组件
func (s *Server) RegisterComponent(lifecycle lifecycle.Lifecycle) {

}

func (s *Server) Run() {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	lifecycle.NewLifecycleManager(s.cfg).Start(ctx, wg)
	log.Info("server start")
	signals := make(chan os.Signal, 1)
	signalList := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGINT,
	}
	signal.Notify(signals, signalList...)
	<-signals
	cancel()
	wg.Wait()
	lifecycle.LifecycleMgr.Stop()
	log.Info("server shutting down")
}
