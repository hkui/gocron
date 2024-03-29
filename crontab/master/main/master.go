package main

import (
	"flag"
	"fmt"
	"gocron/crontab/common"
	"gocron/crontab/master"
	"time"
)

var (
	confFile string
)

func initArgs() {
	// master -config ./master.json
	flag.StringVar(&confFile, "config", "conf/master.json", "指定master.json")
	flag.Parse()
}
func main() {
	var (
		err error
	)
	//初始化线程
	common.InitEnv()
	initArgs()
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	//对etcd的curd操作
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}
	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}
	if err = master.InitWorkerMgr(); err != nil {
		goto ERR
	}
	if err = master.InitUserMgr(); err != nil {
		goto ERR
	}

	//启动Api HTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	fmt.Println("Server listen on ", master.G_config.ApiPort)
	for {
		time.Sleep(1 * time.Second)
	}

	return

ERR:
	fmt.Println(err)

}
