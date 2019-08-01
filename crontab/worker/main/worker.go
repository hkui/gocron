package main

import (
	"crontab/worker"
	"flag"
	"fmt"
	"runtime"
	"time"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
var (
	confFile string
)

func initArgs() {
	// worker -config ./worker.json
	flag.StringVar(&confFile, "config", "./worker.json", "指定worker.josn")
	flag.Parse()
}
func main() {
	var (
		err error
	)
	initArgs()
	//初始化线程
	initEnv()
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}
	//启动执行器
	if err=worker.InitExecutor();err!=nil{
		goto ERR
	}

	if err=worker.InitScheduler();err!=nil{
		goto ERR
	}
	if err=worker.InitLogSink();err!=nil{
		goto ERR
	}

	//初始化任务管理器
	if err=worker.InitJobMgr();err!=nil{
		goto ERR
	}

	err=worker.G_jobMgr.WatchJobs()
	if err!=nil{
		goto ERR
	}
	worker.G_jobMgr.WatchKiller()

	for {
		time.Sleep(1 * time.Second)
	}
	return

ERR:
	fmt.Println(err)

}
