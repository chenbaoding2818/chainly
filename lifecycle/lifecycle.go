package lifecycle

import (
	"sort"
	"sync"

	"github.com/chenbaoding2818/chainly/config"
)

const (
	HighPriority   = 10000000
	NormalPriority = 5000000
	LowPriority    = 0
)

var (
	LifecycleOnce sync.Once
	LifecycleMgr  *LifecycleManager
)

type Lifecycle interface {
	Start(*config.Config)

	Priority() uint32

	Stop()
}

// LifecycleManager 组件生命周期管理器
type LifecycleManager struct {
	lifecycles []Lifecycle
	cfg        *config.Config
}

// Start 启动顺序，从高优先级到低优先级
func (lm *LifecycleManager) Start() {
	if len(lm.lifecycles) == 0 {
		// warning log here no any component registered
		return
	}
	sort.SliceStable(lm.lifecycles, func(i, j int) bool {
		return lm.lifecycles[i].Priority() > lm.lifecycles[j].Priority()
	})
	for _, lifecycle := range lm.lifecycles {
		lifecycle.Start(lm.cfg)
	}
}

// Stop 停止顺序，跟启动顺序相反
func (lm *LifecycleManager) Stop() {
	sort.SliceStable(lm.lifecycles, func(i, j int) bool {
		return lm.lifecycles[i].Priority() < lm.lifecycles[j].Priority()
	})
	for _, lifecycle := range lm.lifecycles {
		lifecycle.Stop()
	}
}

func (lm *LifecycleManager) AddLifecycle(lifecycle Lifecycle) {
	lm.lifecycles = append(lm.lifecycles, lifecycle)
}

func NewLifecycleManager(cfg *config.Config) *LifecycleManager {
	LifecycleOnce.Do(func() {
		if LifecycleMgr == nil {
			LifecycleMgr = &LifecycleManager{
				lifecycles: make([]Lifecycle, 0),
				cfg:        cfg,
			}
		}
	})
	return LifecycleMgr
}
