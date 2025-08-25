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

type ActorManager struct {
	// key: actor id[player id or other id], value: actor
	actors  sync.Map
	handler iface.IMessageHandler
	// TODO: 使用分段锁实现 提高并发能力
	lock sync.Mutex
}

func (am *ActorManager) AddActor(actorId string, conn iface.IConnection) {
	// 判断是否有旧的actor 有的话要删除旧的（异端踢号）保证玩家只能在一个端登陆

	am.actors.Store(actorId, &Actor{
		Id:      actorId,
		Conn:    conn,
		handler: am.handler,
		mailbox: make(chan []byte),
		quitCh:  make(chan struct{}),
	})
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
