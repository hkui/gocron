package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"gocron/base/etcd/client"
	"time"
)

func main() {
	var (
		Connectclient *clientv3.Client
		err           error
	)
	if Connectclient, err = client.GetClient(); err != nil {
		fmt.Println(err)
	}

	lease := clientv3.NewLease(Connectclient)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	//申请一个10s租约
	leaseGrantResp, err := lease.Grant(ctx, 20)

	cancel()
	if err != nil {
		fmt.Println(err)
		return
	}
	//租约Id
	leaseId := leaseGrantResp.ID

	ctx, cancelFun := context.WithCancel(context.Background())
	//5秒后停止续租
	time.AfterFunc(1005*time.Second, func() {
		cancelFun()
	})

	//自动续租
	keepRespChan, err := lease.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}
	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约已失效")
					goto END
				} else {
					fmt.Println("收到续租应答", keepResp.ID, time.Now())
				}
			}
		}
	END:
	}()

	kv := clientv3.NewKV(Connectclient)

	//put一个kv,让它与租约关联起来,从而实现10s
	key := "/cron/lock/job1"
	putResp, err := kv.Put(context.TODO(), key, "job1", clientv3.WithLease(leaseId))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("写入成功:%d\n", putResp.Header.Revision)

	for {
		getResp, err := kv.Get(context.TODO(), key)
		if err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv 过期了")
			break
		}
		fmt.Println("还没过期", time.Now())
		time.Sleep(1 * time.Second)
	}

}
