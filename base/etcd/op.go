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
		opResp clientv3.OpResponse
	)
	if Connectclient,err=client.GetClient();err!=nil{
		fmt.Println(err)
		return
	}
	kv:=clientv3.NewKV(Connectclient)
	key:="/jobs/job1"

	putOp:=clientv3.OpPut(key,"123")

	//执行op
	if opResp,err=kv.Do(context.TODO(),putOp);err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("写入Revision:",opResp.Put().Header.Revision)

	getOp:=clientv3.OpGet(key)

	//执行Op
	if opResp,err=kv.Do(context.TODO(),getOp);err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(opResp.Get().Kvs)














}
