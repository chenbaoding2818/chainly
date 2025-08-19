package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
	// "github.com/chenbaoding2818/chainly/lifecycle/log"
)

type Server struct {
	cfg *config.Config
}

type Option func(*config.Config)

func NewServer(cfg *config.Config, opts ...Option) *Server {
	for _, opt := range opts {
		opt(cfg)
	}
	return &Server{
		cfg: cfg,
	}
}

func NewDefaultServer(opts ...Option) *Server {
	return NewServer(config.NewDefaultConfig(), opts...)
}

func WithNetConfig(netCfg *config.Net) Option {
	return func(c *config.Config) {
		c.NetCfg = netCfg
	}
}

// WithCustomConfig 自定义配置
func WithCustomConfig(m map[string]interface{}) Option {
	return func(cfg *config.Config) {
		for k, v := range m {
			cfg.CustomCfgFields[k] = v
		}
	}
}

// RegisterDefaultComponents 注册默认服务组件
func (s *Server) RegisterDefaultComponents() {
	lifecycle.LifecycleMgr.AddLifecycle(nil)
}

// RegisterComponent 注册服务组件
func (s *Server) RegisterComponent(lifecycle lifecycle.Lifecycle) {

}

func (s *Server) Run() {
	lifecycle.NewLifecycleManager(s.cfg).Start()
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
