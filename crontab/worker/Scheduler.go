package worker

import (
	"crontab/common"
	"fmt"
	"time"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent  //etcd任务事件队列
	jobPlanTable map[string]*common.JobSchedulePlan  //任务调度计划表
	jobExecuteInfoTable map[string]*common.JobExecuteInfo
}
var(
	G_scheduler *Scheduler
)

func InitScheduler() (err error) {
	G_scheduler=&Scheduler{
		jobEventChan:make(chan *common.JobEvent,1000),
		jobPlanTable:make(map[string]*common.JobSchedulePlan),
		jobExecuteInfoTable:make(map[string]*common.JobExecuteInfo),
	}
	go G_scheduler.scheduleLoop()
	return
}
func (scheduler *Scheduler)scheduleLoop()  {
	var (
		jobEvent *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)
	scheduleAfter=scheduler.TrySchedule()
	//调度的延迟定时器
	scheduleTimer=time.NewTimer(scheduleAfter)
	for{
		select {
			case jobEvent=<-scheduler.jobEventChan:
				scheduler.handleJobEvent(jobEvent)
			case <-scheduleTimer.C://最近的任务到期了
		}
		fmt.Printf("try schedule %+v\n",jobEvent)
		scheduleAfter=scheduler.TrySchedule()
		//重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}
//尝试执行任务
func (scheduler *Scheduler)TryStartJob(plan *common.JobSchedulePlan)  {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
	)
	if jobExecuteInfo ,jobExecuting=scheduler.jobExecuteInfoTable[plan.Job.Name];jobExecuting{
		//跳过执行
		return
	}
	jobExecuteInfo=common.BuildJobExecuteInfo(plan)
	scheduler.jobExecuteInfoTable[plan.Job.Name]=jobExecuteInfo
	//todo 执行

}


//处理任务事件
func (scheduler *Scheduler)handleJobEvent(jobEvent *common.JobEvent)  {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		err error
		jobExisted bool
	)
	fmt.Println("handleJobEvent",jobEvent.Job.Name,jobEvent.EventType)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if jobSchedulePlan,err=common.BuildJobSchedulePlan(jobEvent.Job);err!=nil{
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name]=jobSchedulePlan

	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan,jobExisted=scheduler.jobPlanTable[jobEvent.Job.Name];jobExisted{
			delete(scheduler.jobPlanTable,jobEvent.Job.Name)
		}

	}
}
func (scheduler *Scheduler)PushJobEvent(jobEvent *common.JobEvent)  {
	scheduler.jobEventChan<-jobEvent
}
//重新计算任务调度状态

func (scheduler *Scheduler)TrySchedule()(scheduleAfter time.Duration)  {
	var (
		jobPlan *common.JobSchedulePlan
		now time.Time
		nearTime *time.Time
	)
	//如果任务表为空，随便sleep一个时间
	if len(scheduler.jobPlanTable)==0{
		scheduleAfter=1*time.Second
		return
	}
	now=time.Now()
	//遍历所有任务
	for _,jobPlan=range scheduler.jobPlanTable{
		if jobPlan.NextTime.Before(now)||jobPlan.NextTime.Equal(now){
			//todo 尝试执行任务
			fmt.Println("执行",jobPlan.Job.Name,time.Now().Format("2006-01-02 15:04:05"))
			jobPlan.NextTime=jobPlan.Expr.Next(now)
			//统计最近一个要过期的任务事件
		}
		if nearTime==nil||jobPlan.NextTime.Before(*nearTime){
				nearTime=&jobPlan.NextTime
		}
	}

	//下次调度时间间隔
	scheduleAfter=(*nearTime).Sub(now)
	return

}
