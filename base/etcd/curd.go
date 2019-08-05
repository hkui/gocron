package main

import (
	"base/etcd/client"
	"context"
	"crontab/common"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (

		Connectclient *clientv3.Client
		err error
		jobList common.JobList
		jobListOne common.JobListOne

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
		clientv3.WithLimit(3),
		clientv3.WithMaxModRev(0),
		clientv3.WithMinModRev(10338),

	)


	cancel()

	if err!=nil{
		fmt.Println(err)
		return
	}

	//读
	fmt.Println(Resp.Count,Resp.More,Resp.Header)
	for _,kvPair:=range Resp.Kvs{
		if err=json.Unmarshal(kvPair.Value,&jobListOne);err==nil{
			jobListOne.ModRevision=kvPair.ModRevision
			jobList=append(jobList,jobListOne)
			fmt.Println(kvPair)
		}else{
			err=nil
		}
	}
	//sort.Sort(common.JobList(jobList))
	//fmt.Printf("%++v\n",jobList)





}
