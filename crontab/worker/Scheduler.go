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
	jobResultChan chan *common.JobExecuteResult //任务结果队列
}
var(
	G_scheduler *Scheduler
)

func InitScheduler() (err error) {
	G_scheduler=&Scheduler{
		jobEventChan:make(chan *common.JobEvent,1000),
		jobPlanTable:make(map[string]*common.JobSchedulePlan),
		jobExecuteInfoTable:make(map[string]*common.JobExecuteInfo),
		jobResultChan:make(chan *common.JobExecuteResult),
	}
	go G_scheduler.scheduleLoop()
	return
}
func (scheduler *Scheduler)scheduleLoop()  {
	var (
		jobEvent *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult *common.JobExecuteResult
	)
	scheduleAfter=scheduler.TrySchedule()
	//调度的延迟定时器
	scheduleTimer=time.NewTimer(scheduleAfter)
	for{
		select {
			case jobEvent=<-scheduler.jobEventChan:
				scheduler.handleJobEvent(jobEvent)  // 计划调度表的curd
			case <-scheduleTimer.C://最近的任务到期了
				fmt.Println("timer",scheduleAfter)
			case jobResult=<-scheduler.jobResultChan:
				scheduler.HandleJobResult(jobResult)

		}
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
	//任务正在执行中还未结束，跳过本次调度
	if jobExecuteInfo ,jobExecuting=scheduler.jobExecuteInfoTable[plan.Job.Name];jobExecuting{
		//跳过执行
		fmt.Println("正在执行，还没有退出,跳过",jobExecuteInfo.Job.Name)
		return
	}
	//构建执行状态信息
	jobExecuteInfo=common.BuildJobExecuteInfo(plan)
	//保存执行状态
	scheduler.jobExecuteInfoTable[plan.Job.Name]=jobExecuteInfo
	//执行任务
	fmt.Println("执行任务:",jobExecuteInfo.Job.Name,jobExecuteInfo.PlanTime,jobExecuteInfo.RealTime)
	G_executor.ExecuteJob(jobExecuteInfo)

}


//处理任务事件
func (scheduler *Scheduler)handleJobEvent(jobEvent *common.JobEvent)  {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		err error
		jobExisted bool
		jobExecuteinfo *common.JobExecuteInfo
		jobExecuting bool

	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if jobSchedulePlan,err=common.BuildJobSchedulePlan(jobEvent.Job);err!=nil{
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name]=jobSchedulePlan

	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan,jobExisted=scheduler.jobPlanTable[jobEvent.Job.Name];jobExisted{
			fmt.Println("删除了",jobEvent.Job.Name)
			delete(scheduler.jobPlanTable,jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL:
		if jobExecuteinfo,jobExecuting=scheduler.jobExecuteInfoTable[jobEvent.Job.Name];jobExecuting{
			jobExecuteinfo.CancelFunc()
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
			fmt.Println("将要执行",jobPlan.Job.Name,time.Now())
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime=jobPlan.Expr.Next(now)
			//统计最近一个要过期的任务事件
		}
		//查找最近要过期的
		if nearTime==nil||jobPlan.NextTime.Before(*nearTime){
				nearTime=&jobPlan.NextTime
		}
	}

	//下次调度时间间隔
	scheduleAfter=(*nearTime).Sub(now)
	return

}
//执行完的结果队列
func (scheduler *Scheduler)PushJobResult(result *common.JobExecuteResult)  {
	scheduler.jobResultChan<-result
}
//处理任务结果
func (scheduler *Scheduler)HandleJobResult(result *common.JobExecuteResult)  {
	//删除执行状态
	delete(scheduler.jobExecuteInfoTable,result.ExecuteInfo.Job.Name)
	fmt.Printf("完成:job=%s,output=%s,err=%v",result.ExecuteInfo.Job.Name,string(result.Output),result.Err)
	var (
		jobLog *common.JobLog
	)
	if result.Err!=common.ERR_LOCK_ALREADY_REQUIRED{
		jobLog=&common.JobLog{
			JobName:result.ExecuteInfo.Job.Name,
			Command:result.ExecuteInfo.Job.Command,
			OutPut:string(result.Output),
			PlanTime:result.ExecuteInfo.PlanTime.UnixNano()/1000/1000,
			ScheduleTime:result.StartTime.UnixNano()/1000/1000,
			StartTime:result.StartTime.UnixNano()/1000/1000,
			EndTime:result.EndTime.UnixNano()/1000/1000,

		}
		if result.Err!=nil{
			jobLog.Err=result.Err.Error()
		}else{
			jobLog.Err=""
		}
		//todo 存储到mongodb
	}

}
