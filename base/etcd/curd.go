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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//写
	//Resp,err:=kv.Put(ctx,"/cron/jobs/job2","job2",clientv3.WithPrevKV())
	//读某一个
	//Resp,err:=kv.Get(ctx,"/cron/jobs/job2")

	//读某个前缀的

	Resp,err:=kv.Get(ctx,"/cron/jobs/",
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByModRevision,clientv3.SortAscend),
		clientv3.WithLimit(6),


	)

	cancel()

	if err!=nil{
		fmt.Println(err)
		return
	}

	//读
	fmt.Printf("%+v\n",Resp.Kvs)
	for k,v:=range Resp.Kvs{
		fmt.Println(k,v)
	}





}
