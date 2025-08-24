package core

import "github.com/chenbaoding2818/chainly/config"

type Option func(*config.Config)

// WithCustomConfig 自定义配置
func WithCustomConfig(m map[string]interface{}) Option {
	return func(cfg *config.Config) {
		for k, v := range m {
			cfg.CustomCfgFields[k] = v
		}
	}
}

// WithNetConfig 网络配置定义
func WithNetConfig(netCfg *config.Net) Option {
	return func(c *config.Config) {
		c.NetCfg = netCfg
	}
}

// WithLogConfig 日志配置定义
func WithLogConfig(logCfg *config.Log) {

}

// WithLogConfig 存储配置定义
func WithStorageConfig(logCfg *config.Storage) {

}
