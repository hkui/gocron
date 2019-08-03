package common

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

func GetEtcdClient(EtcdEndpoints []string,EctdDialTimeout int64)(client *clientv3.Client,err error )  {
	var(
		config clientv3.Config
	)
	config=clientv3.Config{
		Endpoints:EtcdEndpoints,//集群地址
		DialTimeout:time.Millisecond*time.Duration(EctdDialTimeout),
	}
	//建立连接
	if client,err=clientv3.New(config);err!=nil{
		return
	}
	return
}
