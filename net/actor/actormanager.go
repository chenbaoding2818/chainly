package actor

import (
	"sync"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
)

var (
	ActorMgr     *ActorManager
	ActorMgrOnce sync.Once
)

// TODO: 使用分段锁实现 提高并发能力
type ActorManager struct {
	// key: actor id[player id or other id], value: actor
	actors  sync.Map
	handler iface.IMessageHandler
	lock    sync.Mutex
}

func (am *ActorManager) AddActor(actorId string, actor *Actor) {
	// 判断是否有旧的actor 有的话要删除旧的（异端踢号）保证玩家只能在一个端登陆
	// TODO:待完善
	if v, ok := am.actors.Load(actorId); ok {
		oldActor := v.(*Actor)
		if oldActor.Conn != nil {
			// 旧的actor有连接 则通知旧的actor 断开连接
			oldActor.Conn.WriteMessage(oldActor.otherClientMsg)
			cancel := oldActor.Conn.GetCancel()
			// 退出旧的连接监听
			cancel()
			oldActor.Stop()
		}
	}
	am.actors.Store(actorId, actor)
}

func (am *ActorManager) RemoveActor(actorId string) {
	am.actors.Delete(actorId)
}

func (am *ActorManager) GetActor(actorId string) *Actor {
	am.lock.Lock()
	defer am.lock.Unlock()
	if actor, ok := am.actors.Load(actorId); ok {
		// 判断连接是否存在
		if actor.(*Actor).Conn != nil {

		}
		return actor.(*Actor)
	}
	// 不存在 则创建新actor
	actor := &Actor{
		Id:      actorId,
		mailbox: make(chan []byte),
		quitCh:  make(chan struct{}),
		handler: am.handler,
	}
	am.actors.Store(actorId, actor)
	// 启动actor
	go actor.start()
	return actor
}

func (am *ActorManager) Run() {
	go func() {
		// TODO: 实现actor清除管理
	}()
}

func NewActorManager(cfg *config.Net) *ActorManager {
	ActorMgrOnce.Do(func() {
		if ActorMgr == nil {
			ActorMgr = &ActorManager{
				actors:  sync.Map{},
				handler: cfg.MessageHandler,
			}
		}
	})
	return ActorMgr
}
