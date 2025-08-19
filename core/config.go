package core

// 服务基础配置文件解析器类型
type ConfigParser uint8

const (
	IniParser ConfigParser = iota + 1
	JsonParser
	YamlParser
)

// Config 定义服务基础配置文件结构
type Config struct {
	// 配置路径，默认是当前目录下的config.ini
	Path string
	// 配置解析器，默认是ini，支持json、yaml
	Parser ConfigParser
	// 高并发自动扩容 放到默认配置中
	concurrencyAutoEnabe bool
}

func NewDefaultConfig() *Config {
	return &Config{
		Path:   "./server.ini",
		Parser: IniParser,
	}
}
