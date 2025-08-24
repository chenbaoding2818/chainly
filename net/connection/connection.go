package connection

import (
	"sync"

	iface "github.com/chenbaoding2818/chainly/interface"
)

var (
	connMgr     *ConnectionManager
	connMgrOnce sync.Once
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	// ws连接map key: 连接的地址， value: 连接对象
	Conns sync.Map // 连接池
}

func (cm *ConnectionManager) AddConn(ptr uintptr, conn iface.IConnection) {
	cm.Conns.Store(ptr, conn)
}

func (cm *ConnectionManager) DelConn(conn iface.IConnection) {

}

func (cm *ConnectionManager) GetConn(addr string) iface.IConnection {
	return nil
}

// Len 连接数量
func (cm *ConnectionManager) Len() int {
	var count int
	cm.Conns.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	return count
}

func NewConnManager() *ConnectionManager {
	connMgrOnce.Do(func() {
		if connMgr == nil {
			connMgr = &ConnectionManager{}
		}
	})
	return connMgr
}
