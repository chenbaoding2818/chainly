package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chenbaoding2818/chainly/lifecycle"
	// "github.com/chenbaoding2818/chainly/lifecycle/log"
)

type Server struct {
	ConfigPath *Config
}

type Option func(*Config)

func NewServer(opts ...Option) *Server {
	cfg := NewDefaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}
	return &Server{}
}

// RegisterDefaultComponents 注册默认服务组件
func (s *Server) RegisterDefaultComponents() {
	lifecycle.LifecycleMgr.AddLifecycle(nil)
}

// RegisterComponent 注册服务组件
func (s *Server) RegisterComponent(lifecycle lifecycle.Lifecycle) {

}

func (s *Server) Run() {
	lifecycle.LifecycleMgr.Start()
	// log.Info("server start")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-signals:
		lifecycle.LifecycleMgr.Stop()
	}
	// log.Info("server shutting down")
}
