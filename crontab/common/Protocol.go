package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}
//调度计划
type JobSchedulePlan struct {
	Job *Job
	Expr *cronexpr.Expression
	NextTime time.Time
}
//任务执行状态
type JobExecuteInfo struct {
	Job *Job
	PlanTime time.Time  //理论上的调度时间
	RealTime time.Time
}
//HTTP接口应答
type Response struct {
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
type JobEvent struct {
	EventType int
	Job *Job
}

//应答方法

func BuildResponse(errno int,msg string,data interface{})(resp []byte,err error)  {
	var (
		response Response
	)
	response.Errno=errno
	response.Msg=msg
	response.Data=data
	resp,err=json.Marshal(response)
	return

}
//反序列化,josn转为job
func UnpackJob(value []byte)(ret *Job,err error)  {
	var (
		job *Job
	)
	job=&Job{}
	if err=json.Unmarshal(value,job);err!=nil{
		return
	}
	ret=job
	return
}
// 从/cron/jobs/job1  得到job1
func ExtraJobName(jobKey string)string  {
	return strings.TrimPrefix(jobKey,JOB_KILLER_DIR)
}
func BuildJobEvent(eventType int,job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType:eventType,
		Job:job,
	}
}

//构造任务执行计划
func BuildJobSchedulePlan(job *Job)(jobSchedulePlan *JobSchedulePlan,err error)  {
	var (
		expr *cronexpr.Expression
	)
	if expr,err=cronexpr.Parse(job.CronExpr);err!=nil{
		return
	}
	//生成任务调度计划对象
	jobSchedulePlan=&JobSchedulePlan{
		Job:job,
		Expr:expr,
		NextTime:expr.Next(time.Now()),
	}
	return
}

func BuildJobExecuteInfo( paln  *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo=&JobExecuteInfo{
		Job:paln.Job,
		PlanTime:paln.NextTime,
		RealTime:time.Now(),
	}
	return
}

//检查cron表达式并返回下10次的执行时间
func CheckCronExpr(mycronExp string)(nexts []string,err error)  {
	var(

		exp *cronexpr.Expression
		now time.Time
		nextTime time.Time
	)
	exp,err=cronexpr.Parse(mycronExp)
	if err!=nil{
		return
	}
	now=time.Now()
	nextTime=exp.Next(now)

	for i:=0;i<10;i++{
		if i==0{
			nextTime=exp.Next(now)
		}else{
			nextTime=exp.Next(nextTime)
		}
		nexts= append(nexts, nextTime.Format("2006-01-02 15:04:05"))
	}
	return
}
