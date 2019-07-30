package main

import (
	"crontab/master"
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
	// master -config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "指定master.josn")
	flag.Parse()
}
func main() {

	var (
		err error
	)
	initArgs()
	//初始化线程
	initEnv()
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}
	//启动Api HTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	for {
		time.Sleep(1 * time.Second)
	}
	return

ERR:
	fmt.Println(err)

}
