package connection

import "sync"

// ConnectionManager 连接管理器
type ConnectionManager struct {
	Conns sync.Map // 连接池
}
