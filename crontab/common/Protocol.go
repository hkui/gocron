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
