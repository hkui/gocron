package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

//调度多个

type CronJob struct {
	expr *cronexpr.Expression
	nextTime time.Time
}


func main() {
	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		now time.Time
		scheduleTable map[string]*CronJob
	)
	now=time.Now()
	scheduleTable=make(map[string]*CronJob)

	expr=cronexpr.MustParse("*/2 * * * * * *")

	cronJob=&CronJob{
		expr:expr,
		nextTime:expr.Next(now),
	}
	scheduleTable["job1"]=cronJob

	expr=cronexpr.MustParse("*/2 * * * * * *")

	cronJob=&CronJob{
		expr:expr,
		nextTime:expr.Next(now),
	}
	scheduleTable["job2"]=cronJob

	go func() {
		var(
			jobName string
			cronJob *CronJob
			now time.Time
		)

		//定时检查任务调度表
		for   {
			now=time.Now()

			for jobName,cronJob=range scheduleTable{
				//判断是否过期
				if cronJob.nextTime.Before(now)||cronJob.nextTime.Equal(now){
					//启动协程执行任务
					go func(jobName string) {
						fmt.Println(jobName,"execed")
					}(jobName)

					//计算下次执行时间
					cronJob.nextTime=cronJob.expr.Next(now)
					fmt.Printf("%s 在 %s下次执行\n",jobName,cronJob.nextTime)
				}
			}
			select {
				case <-time.NewTicker(100*time.Millisecond).C:  //将在100毫秒后可读返回
			}
		}
	}()
	time.Sleep(100*time.Second)


}
