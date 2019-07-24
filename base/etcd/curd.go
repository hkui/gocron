package main

import (
	"base/etcd/client"
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (

		Connectclient *clientv3.Client
		err error

	)
	if Connectclient,err=client.GetClient();err!=nil{
		fmt.Println(err);
	}


	kv:=clientv3.NewKV(Connectclient)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	//写
	//Resp,err:=kv.Put(ctx,"/cron/jobs/job2","job2",clientv3.WithPrevKV())
	//读某一个
	//Resp,err:=kv.Get(ctx,"/cron/jobs/job2")

	//读某个前缀的
	Resp,err:=kv.Get(ctx,"/cron/jobs/",clientv3.WithPrefix())

	cancel()

	if err!=nil{
		fmt.Println(err)
		return
	}
	//读
	fmt.Printf("%+v\n",Resp.Kvs)





}
