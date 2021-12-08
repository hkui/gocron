package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"gocron/base/etcd/client"
	"strings"
	"time"
)

func main() {
	var (
		Connectclient *clientv3.Client
		err           error
		key           string
		i             int
		WatchChan     clientv3.WatchChan
	)
	key = "hk"

	if Connectclient, err = client.GetClient(); err != nil {
		fmt.Println(err)
	}
	kv := clientv3.NewKV(Connectclient)
	//往里面写数据,改数据删数据
	go func() {
		for i = 0; i < 10; i++ {
			val := strings.Join([]string{"val", string(i)}, "-")
			if i%2 == 0 {
				kv.Put(context.TODO(), key, val)
			} else {
				kv.Delete(context.TODO(), key)
			}

			time.Sleep(time.Second)
		}

	}()

	WatchChan = Connectclient.Watch(context.TODO(), key)

	for v := range WatchChan {

		for _, vv := range v.Events {
			fmt.Println(vv, vv.Type, string(vv.Kv.Key))
		}
	}

}
