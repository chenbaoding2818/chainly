package cron

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	CronComponent      *CronManager
	TimerComponentOnce sync.Once
)

type TimerTask struct {
	Spec string
	Cmd  func()
}

type CronManager struct {
	Cron          *cron.Cron
	TimerTaskList []TimerTask // 定时任务列表
}

func (cc *CronManager) RegisterTimers(taskList []TimerTask) {
	cc.TimerTaskList = append(cc.TimerTaskList, taskList...)
}

func (cc *CronManager) RegisterTimer(task TimerTask) {
	cc.TimerTaskList = append(cc.TimerTaskList, task)
}

func NewCronComponent() *CronManager {
	TimerComponentOnce.Do(func() {
		if CronComponent == nil {
			CronComponent = &CronManager{
				Cron:          cron.New(cron.WithSeconds()),
				TimerTaskList: make([]TimerTask, 0),
			}
		}
	})
	return CronComponent
}
