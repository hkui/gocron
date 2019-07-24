package client

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func GetClient()(*clientv3.Client,error){
	var (
		config clientv3.Config
		Connectclient *clientv3.Client
		err error

	)
	config=clientv3.Config{
		Endpoints:[]string{"39.100.78.46:2379"},
		DialTimeout:3*time.Second,
	}

	if Connectclient,err=clientv3.New(config);err!=nil{
		fmt.Println(err)
		return nil,err
	}
	return Connectclient,nil;
}