package main

import (
	"base/etcd/client"
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
)

func main() {
	var (

		Connectclient *clientv3.Client
		err error
		putResp *clientv3.PutResponse

	)
	if Connectclient,err=client.GetClient();err!=nil{
		fmt.Println(err);
	}
	


	kv:=clientv3.NewKV(Connectclient)
	if putResp,err=kv.Put(context.TODO(),"/cron/jobs/job1","job1");err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(putResp)



}
