package main

import (
	"base/etcd/client"
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//lease实现锁自动过期
//op操作
//txn事务
func main() {
	var (
		Connectclient *clientv3.Client
		err error
		txnResp *clientv3.TxnResponse

	)
	if Connectclient,err=client.GetClient();err!=nil{
		fmt.Println(err)
		return
	}
	/**
	1 上锁(创建租约,自动续租,拿着租约去抢占一个key)
	2 处理业务
	3.释放锁，取消自动续租 释放租约(会立即删除key)

	 */
	key:="/jobs/job1"

	lease:=clientv3.NewLease(Connectclient)
	//申请一个10s租约
	leaseGrantResp,err:=lease.Grant(context.TODO(),10)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//租约Id
	leaseId:=leaseGrantResp.ID
	//创建一个用于取消续租的context

	ctx, cancel := context.WithCancel(context.TODO())

	//确保函数退出后自动续租会立刻停止
	defer  cancel()
	defer lease.Revoke(context.TODO(),leaseId)

	//自动续租
	keepRespChan,err:=lease.KeepAlive(ctx,leaseId)
	if err!=nil{
		fmt.Println(err)
		return
	}
	//处理续租应答的协程
	go func() {
		for{
			select {
				case keepResp:=<-keepRespChan:
					if keepResp ==nil{
						fmt.Println("租约已失效")
						goto  END
					}else{
						fmt.Println("收到续租应答",keepResp.ID)
					}
			}
		}
		END:
	}()
	// if 不存在Key,then设置它else 抢锁失败


	kv:=clientv3.NewKV(Connectclient)
	//kv.Delete(context.TODO(),key)
	//return
	//创建事务

	txn:=kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(key),"=",0)).
		Then(clientv3.OpPut(key,"",clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(key)) //抢锁失败

		if txnResp,err=txn.Commit();err!=nil{
			fmt.Println(err)
			return
		}
		//判断是否抢到了锁
		if !txnResp.Succeeded{
			fmt.Println("锁被占用:",string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		}
		//处理业务
		fmt.Println("业务处理")

		time.Sleep(5*time.Second)

		//释放锁（取消自动续租,释放续租）
		//defer会把租约释放掉,关联的kv就被删除

}
