package connection

import (
	"sync"

	iface "github.com/chenbaoding2818/chainly/interface"
	"github.com/chenbaoding2818/chainly/net/actor"
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

func (cm *ConnectionManager) AddConn(actor *actor.Actor) {
	cm.Conns.Store(actor.Conn.GetConnPtr(), actor)
}

func (cm *ConnectionManager) DelConn(conn iface.IConnection) {

}

func (cm *ConnectionManager) GetConn(addr string) iface.IConnection {
	return nil
}

// Len 连接数量 TODO: 目前使用遍历的手段统计连接数，后期可以考虑维护一个计数器，但是复杂度会增加
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
