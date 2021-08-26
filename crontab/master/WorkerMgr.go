package master

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"gocron/crontab/common"
)

// /cron/workers/
type WorkerMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var (
	G_workerMgr *WorkerMgr
)
func InitWorkerMgr() (err error) {
	var (
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)


	// 建立连接
	if client, err = common.GetEtcdClient(G_config.EtcdEndpoints,G_config.EctdDialTimeout); err != nil {
		return
	}

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_workerMgr = &WorkerMgr{
		client :client,
		kv: kv,
		lease: lease,
	}
	return
}

// 获取在线worker列表
func (workerMgr *WorkerMgr) ListWorkers() (workerArr []string, err error) {
	var (
		getResp *clientv3.GetResponse
		kv *mvccpb.KeyValue
		workerIP string
	)

	// 初始化数组
	workerArr = make([]string, 0)

	// 获取目录下所有Kv
	if getResp, err = workerMgr.kv.Get(
		context.TODO(),
		common.JOB_WORKER_DIR,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByModRevision,clientv3.SortDescend),
		); err != nil {
		return
	}

	// 解析每个节点的IP
	for _, kv = range getResp.Kvs {
		// kv.Key : /cron/workers/192.168.2.1
		workerIP = common.ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, workerIP)
	}
	return
}


