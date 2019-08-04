package main

import (
	"crontab/common"
	"crontab/master"
	"flag"
	"fmt"
	"time"
)


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
	common.InitEnv()
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}
	//对etcd的curd操作
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}
	if err=master.InitLogMgr();err!=nil{
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
