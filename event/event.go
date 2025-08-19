package event

import (
	"runtime/debug"
	"sort"
)

type EventType uint16

// 事件类型
const (
	// 创建角色事件
	CreateRoleEvent EventType = iota
	// 登录事件
	LoginEvent
	// 登出事件 内部事件 经由此事件转发到channel再转发到ZeroEvent保证线性处理，没有数据竞争
	LogoutInternalEvent
	// 登出事件
	LogoutEvent
	// 心跳事件
	HeartBeatEvent
	// 00:00事件 每个玩家用于刷新凡是00:00需要刷新的任务
	ZeroEvent
	// 00:00事件 内部事件 经由此事件转发到channel再转发到ZeroEvent保证线性处理，没有数据竞争
	ZeroEventInternal
	// 充值事件
	RechargeEvent
)

const (
	HighPriority   = 10000000
	NormalPriority = 5000000
	LowPriority    = 0
)

type Event struct {
	Type   EventType
	Player any
	Args   any
}

type EventListener struct {
	Priority uint32
	OnEvent  func(event *Event)
}

type eventManager struct {
	events map[EventType][]*EventListener
}

var (
	emInstance *eventManager
)

func (em *eventManager) registerEvent(eventType EventType, listener *EventListener) {
	list := em.events[eventType]
	if list == nil {
		list = make([]*EventListener, 0)
	}
	list = append(list, listener)
	em.events[eventType] = list
}

func recoverPanic(listener *EventListener, event *Event) {
	defer func() {
		if err := recover(); err != nil {
			// log.Error(fmt.Sprint(err))
			debug.PrintStack()
		}
	}()
	listener.OnEvent(event)
}

func (em *eventManager) fireEvent(eventType EventType, event *Event, syncAction bool) {
	list := em.events[eventType]
	if list == nil {
		return
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Priority > list[j].Priority
	})
	for _, listener := range list {
		listen := listener
		if syncAction {
			recoverPanic(listen, event)
		} else {
			// TODO: 注意使用协程池
			go recoverPanic(listen, event)
			// worker.AddTask(func() {
			// 	recoverPanic(listen, event)
			// })
		}
	}
}

// RegisterEvent 同步事件
func RegisterEvent(eventType EventType, listener *EventListener) {
	if emInstance == nil {
		emInstance = &eventManager{
			events: map[EventType][]*EventListener{},
		}
	}
	emInstance.registerEvent(eventType, listener)
}

func FireEvent(eventType EventType, event *Event) {
	emInstance.fireEvent(eventType, event, true)
}

// FireEventAsync 触发异步事件
func FireEventAsync(eventType EventType, event *Event) {
	emInstance.fireEvent(eventType, event, false)
}
