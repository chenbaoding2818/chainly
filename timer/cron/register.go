package cron

import (
	"github.com/chenbaoding2818/chainly/lifecycle"
)

func (cc *CronComponent) Start() {
	for _, task := range cc.TimerTaskList {
		if _, err := cc.Cron.AddFunc(task.Spec, task.Cmd); err != nil {
			// 定时器注册失败
			panic("The crontab was not added successfully.")
		}
	}
	cc.Cron.Start()
}

func (cc *CronComponent) Priority() int32 {
	// 保证定时器执行优先级最低
	return lifecycle.LowPriority - 100
}

func (cc *CronComponent) Stop() {

}
