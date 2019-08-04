package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"runtime"
	"strings"
	"time"
)
func InitEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
type Job struct {
	Name string `json:"name" bson:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}
type JobListOne struct {
	Job
	ModRevision int64
}
type JobList []JobListOne

func (s JobList) Len() int { return len(s) }

func (s JobList) Swap(i, j int){ s[i], s[j] = s[j], s[i] }

func (s JobList) Less(i, j int) bool { return s[i].ModRevision > s[j].ModRevision }




type JobListsRes struct {
	Lists []JobListOne `json:"lists"`
	Sum int64		`json:"sum"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
	NextModRevision int64 `json:"next_mod_revision"`
	PrevModRevision int64 `json:"prev_mod_revision"`
	

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
	CancelCtx context.Context
	CancelFunc context.CancelFunc
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
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output []byte  //脚本输出
	Err error	//脚本错误原因
	StartTime time.Time //启动时间
	EndTime time.Time  //结束时间

}
// 任务执行日志
type JobLog struct {
	JobName string `json:"jobName" bson:"jobName"` // 任务名字
	Command string `json:"command" bson:"command"` // 脚本命令
	Err string `json:"err" bson:"err"` // 错误原因
	Output string `json:"output" bson:"output"`	// 脚本输出
	PlanTime int64 `json:"planTime" bson:"planTime"` // 计划开始时间
	ScheduleTime int64 `json:"scheduleTime" bson:"scheduleTime"` // 实际调度时间
	StartTime int64 `json:"startTime" bson:"startTime"` // 任务执行开始时间
	EndTime int64 `json:"endTime" bson:"endTime"` // 任务执行结束时间
}
type JobLogShow struct {
	JobName string  `json:"jobName" bson:"jobName"`
	Command string  `json:"command" bson:"command"`
	Err string  `json:"err" bson:"err"`
	Output string  `json:"output" bson:"output"`
	PlanTime string `json:"planTime" bson:"planTime"`
	ScheduleTime string `json:"scheduleTime" bson:"scheduleTime"`
	StartTime string  `json:"startTime" bson:"startTime"`
	EndTime string   `json:"endTime" bson:"endTime"`
}

func TimeToStr(time2 time.Time)(tstring string)  {
	tstring=time2.Format("06/01/02 15:04:05.000")
	return
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
	return strings.TrimPrefix(jobKey,JOB_SAVE_DIR)
}
func ExtraKillerName(killerKey string)string  {
	return strings.TrimPrefix(killerKey,JOB_KILLER_DIR)
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
	jobExecuteInfo.CancelCtx,jobExecuteInfo.CancelFunc=context.WithCancel(context.TODO())
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
// 日志批次
type LogBatch struct {
	Logs []interface{}	// 多条日志
}

// 任务日志过滤条件
type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

// 任务日志排序规则
type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`	// {startTime: -1}
}


