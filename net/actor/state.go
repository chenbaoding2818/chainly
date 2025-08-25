package actor

import (
	"sync"
	"time"
)

type PlayerState int

const (
	StateOffline PlayerState = iota //
	StateOnline
	StateProcessing // 临时处理状态（如离线操作中）
)

// StateCoordinator 玩家状态协调器
// 区分离线在线操作，对于玩家离线的操作，创建临时的actor处理，
// 如果在进行离线操作的时候玩家上线了，则将请求转移到在线actor处理
// 增加了玩家状态协调器的优缺点
//
//	优点：对于离线在线操作能够简单处理，不需要复杂的状态同步机制
//	缺点：如果玩家越多协调器所占内存就会越大，需要考虑内存占用问题 TODO:可以考虑使用LRU缓存淘汰机制优化
type StateCoordinator struct {
	stateMap sync.Map // playerID -> *PlayerStateInfo
}

type PlayerStateInfo struct {
	State     PlayerState
	Timestamp int64
	Lock      chan struct{} // 细粒度锁
}

func (sc *StateCoordinator) GetState(playerID int64) PlayerState {
	if info, ok := sc.stateMap.Load(playerID); ok {
		return info.(*PlayerStateInfo).State
	}
	return StateOffline
}

func (sc *StateCoordinator) SetState(playerID int64, state PlayerState) {
	info, ok := sc.stateMap.Load(playerID)
	if ok {
		info.(*PlayerStateInfo).State = state
		info.(*PlayerStateInfo).Timestamp = time.Now().Unix()
	} else {
		lock := make(chan struct{}, 1)
		lock <- struct{}{}
		info = &PlayerStateInfo{
			State:     state,
			Timestamp: time.Now().Unix(),
			Lock:      lock,
		}
		sc.stateMap.Store(playerID, info)
	}
}

// 玩家离线处理
func (sc *StateCoordinator) Offline(playerID int64) {
	info, ok := sc.stateMap.Load(playerID)
	if ok {
		info.(*PlayerStateInfo).State = StateOffline
		info.(*PlayerStateInfo).Timestamp = time.Now().Unix()
	} else {
		lock := make(chan struct{}, 1)
		lock <- struct{}{}
		info = &PlayerStateInfo{
			State:     StateOffline,
			Timestamp: time.Now().Unix(),
			Lock:      lock,
		}
		sc.stateMap.Store(playerID, info)
	}
}

// 玩家上线处理
func (sc *StateCoordinator) Online(playerID int64) {
	info, ok := sc.stateMap.Load(playerID)
	if ok {
		info.(*PlayerStateInfo).State = StateOnline
		info.(*PlayerStateInfo).Timestamp = time.Now().Unix()
	} else {
		lock := make(chan struct{}, 1)
		lock <- struct{}{}
		info = &PlayerStateInfo{
			State:     StateOnline,
			Timestamp: time.Now().Unix(),
			Lock:      lock,
		}
		sc.stateMap.Store(playerID, info)
	}
}

// 玩家临时处理状态处理
func (sc *StateCoordinator) Processing(playerID int64) {
	info, ok := sc.stateMap.Load(playerID)
	if ok {
		info.(*PlayerStateInfo).State = StateProcessing
		info.(*PlayerStateInfo).Timestamp = time.Now().Unix()
	} else {
		lock := make(chan struct{}, 1)
		lock <- struct{}{}
		info = &PlayerStateInfo{
			State:     StateProcessing,
			Timestamp: time.Now().Unix(),
			Lock:      lock,
		}
		sc.stateMap.Store(playerID, info)
	}
}
