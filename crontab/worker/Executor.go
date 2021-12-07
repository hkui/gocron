package worker

import (
	"gocron/crontab/common"
	"os/exec"
	"time"
)

type Executor struct {

}
var (
	G_executor *Executor
)
func InitExecutor()(err error){
	G_executor=&Executor{}
	return
}
//执行一个任务
func (executor *Executor)ExecuteJob(jobExecuteInfo *common.JobExecuteInfo){
	var (
		cmd *exec.Cmd
		err error
		output []byte
		result *common.JobExecuteResult
		jobLock *JobLock
	)
	//执行shell命令
	go func() {
		result=&common.JobExecuteResult{
			ExecuteInfo:jobExecuteInfo,
			Output:make([]byte,0),
			StartTime:time.Now(),
		}
		//初始化锁
		jobLock=G_jobMgr.CreateJobLock(jobExecuteInfo.Job.Name)
		err=jobLock.TryLock()
		defer jobLock.Unlock();
		if err!=nil{
			result.Err=err
			result.EndTime=time.Now()
		}else{
			result.StartTime=time.Now()
			cmd=exec.CommandContext(jobExecuteInfo.CancelCtx,"/bin/bash","-c",jobExecuteInfo.Job.Command)
			output,err=cmd.CombinedOutput()
			result.EndTime=time.Now()
			result.Err=err
			result.Output=output
		}
		//执行结束后把结果返回给调度器,调度器把任务从jobExecuteInfoTable 删除
		G_scheduler.PushJobResult(result)
		//省得多台机器 时间差值大 ，都抢到锁了
		time.Sleep(1*time.Second)
		//释放锁

	}()
}

