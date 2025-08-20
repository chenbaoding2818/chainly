package config

import (
	"strings"
)

// 服务基础配置文件解析器类型
type ConfigParser uint8

const (
	IniParser ConfigParser = iota + 1
	JsonParser
	YamlParser
)

type Basic struct {
}

func NewDefaultBasicCfg() *Basic {
	return &Basic{}
}

// Config 定义服务基础配置文件结构
type Config struct {
	// 配置路径，默认是当前目录下的config.ini
	Path string
	// 配置解析器，默认是ini，支持json、yaml
	Parser ConfigParser
	// 关于服务相关的配置
	BasicCfg *Basic
	// 关于net组件相关的配置
	NetCfg *Net
	// 关于存储的配置
	StorageCfg *Storage
	// 关于日志的配置
	LogCfg *Log
	// 自定义配置项 用户可根据自己的项目自定义自己的配置项
	CustomCfgFields map[string]interface{}
}

func NewConfig(path string, parser ConfigParser) *Config {
	if len(path) == 0 {
		panic("config path is empty")
	}

	pathSplitList := strings.Split(path, ".")
	if len(pathSplitList) < 2 {
		panic("config path is invalid")
	}

	switch pathSplitList[len(pathSplitList)-1] {
	case "ini":
		parser = IniParser
	case "json":
		parser = JsonParser
	case "yaml":
		parser = YamlParser
	default:
		panic("config file type is invalid")
	}
	if parser < IniParser || parser > YamlParser {
		panic("config parser is invalid")
	}

	cfg := NewDefaultConfig()
	cfg.Path = path
	cfg.Parser = parser
	// 配置解析(自定义配置项)
	err := cfg.Parse()
	if err != nil {
		// 解析失败 服务禁止启动
		panic(err)
	}
	return cfg
}

func NewDefaultConfig() *Config {
	// 设置默认的值
	return &Config{
		BasicCfg:        NewDefaultBasicCfg(),
		NetCfg:          NewDefaultNetCfg(),
		CustomCfgFields: make(map[string]interface{}),
	}
}
