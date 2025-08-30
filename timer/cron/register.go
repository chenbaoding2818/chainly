package cron

import (
	"context"

	"github.com/chenbaoding2818/chainly/config"
	"github.com/chenbaoding2818/chainly/lifecycle"
)

func (cc *CronManager) Start(ctx context.Context, cfg *config.Config) {
	for _, task := range cc.TimerTaskList {
		if _, err := cc.Cron.AddFunc(task.Spec, task.Cmd); err != nil {
			// 定时器注册失败
			panic("The crontab was not added successfully.")
		}
	}
	// // Debug模式下每30秒执行一次空任务 为了服务器改时间时定时器能正确运行
	// if config.App.Debug {
	// 	cc.CronManager.AddFunc("*/30 * * * * *", func() {})
	// }
	cc.Cron.Start()
}

func (cc *CronManager) Priority() int32 {
	// 保证定时器执行优先级最低
	return lifecycle.LowPriority - 100
}

func (cc *CronManager) Stop() {

}
