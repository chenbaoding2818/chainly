package cron

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	TimerComponent     *CronComponent
	TimerComponentOnce sync.Once
)

type TimerTask struct {
	Spec string
	Cmd  func()
}

type CronComponent struct {
	Cron          *cron.Cron
	TimerTaskList []TimerTask // 定时任务列表
}

func (cc *CronComponent) RegisterTimers(taskList []TimerTask) {
	cc.TimerTaskList = append(cc.TimerTaskList, taskList...)
}

func (cc *CronComponent) RegisterTimer(task TimerTask) {
	cc.TimerTaskList = append(cc.TimerTaskList, task)
}

func NewCronComponent() *CronComponent {
	TimerComponentOnce.Do(func() {
		if TimerComponent == nil {
			TimerComponent = &CronComponent{
				Cron:          cron.New(cron.WithSeconds()),
				TimerTaskList: make([]TimerTask, 0),
			}
		}
	})
	return TimerComponent
}
