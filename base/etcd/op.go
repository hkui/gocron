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
		lease clientv3.Lease
		leaseResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID

		KeepRespChan  <-chan *clientv3.LeaseKeepAliveResponse

		keepResp *clientv3.LeaseKeepAliveResponse


	)
	if Connectclient,err=client.GetClient();err!=nil{
		fmt.Println(err)
	}
	//创建租约
	lease=clientv3.NewLease(Connectclient)
	//创建1个10s的租约
	leaseResp,err=lease.Grant(context.TODO(),10)

	leaseId=leaseResp.ID


	KeepRespChan,err=lease.KeepAlive(context.TODO(),leaseId)

	//put一个带租约的key




	go func() {
		select {
			case keepResp=<-KeepRespChan:
				if keepResp==nil{
					fmt.Println("租约已失效")
					goto END
				}else{
					fmt.Println("续约：",keepResp.ID)
				}
		}
		END:

	}()








}
