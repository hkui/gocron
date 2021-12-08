package common

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func GetEtcdClient(EtcdEndpoints []string, EctdDialTimeout int64) (client *clientv3.Client, err error) {
	var (
		config clientv3.Config
	)
	config = clientv3.Config{
		Endpoints:   EtcdEndpoints, //集群地址
		DialTimeout: time.Millisecond * time.Duration(EctdDialTimeout),
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	return
}
func GetMongoClient(MongodbUri []string) (client *mongo.Client, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err = mongo.Connect(ctx,
		&options.ClientOptions{
			Hosts: MongodbUri,
		},
	)
	if err != nil {
		return
	}
	//检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return
	}
	return
}
