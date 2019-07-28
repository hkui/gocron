package master

import (
	"context"
	"crontab/common"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}
var (
	G_jobMgr *JobMgr
)
//初始化管理器
func InitJobMgr() (err error) {
	var(
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	config=clientv3.Config{
		Endpoints:G_config.EtcdEndpoints,//集群地址
		DialTimeout:time.Millisecond*time.Duration(G_config.EctdDialTimeout),
	}
	//建立连接

	if client,err=clientv3.New(config);err!=nil{
		return
	}
	kv=client.KV
	lease=client.Lease
	G_jobMgr=&JobMgr{
		client:client,
		kv:kv,
		lease:lease,
	}
	return

}
func (JobMgr *JobMgr)SaveJob(job *common.Job) (oldJob *common.Job,err error) {
	var (
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
		oldJobObj common.Job

	)

	jobKey=common.JOB_SAVE_DIR+job.Name
	if jobValue,err=json.Marshal(job);err!=nil{
		return
	}
	//保存到etcd
	putResp,err=JobMgr.kv.Put(context.TODO(),jobKey,string(jobValue),clientv3.WithPrevKV())
	if err!=nil{
		return
	}
	//如果是跟新
	if putResp.PrevKv!=nil{
		if err=json.Unmarshal(putResp.PrevKv.Value,&oldJobObj);err!=nil{
			err=nil
			return
		}
		oldJob=&oldJobObj
	}

	return
}
//删除job
func (JobMgr *JobMgr)DeleteJob(name string) (oldJob *common.Job,err error){
	var (
		jobKey string
		deleteResp *clientv3.DeleteResponse
		oldJobObj common.Job
	)
	jobKey=common.JOB_SAVE_DIR+name
	if deleteResp,err=JobMgr.kv.Delete(context.TODO(),jobKey,clientv3.WithPrevKV());err!=nil{
		return
	}
	if len(deleteResp.PrevKvs)>0{
		if err=json.Unmarshal(deleteResp.PrevKvs[0].Value,&oldJobObj);err!=nil{
			err=nil
			return
		}
		oldJob=&oldJobObj
	}
	return
}
func (JobMgr *JobMgr)JobList()(jobList[]common.Job,err error)  {
	var (
		jobKey string
		getResp *clientv3.GetResponse

		kvPair *mvccpb.KeyValue
		job common.Job
	)
	jobKey=common.JOB_SAVE_DIR
	jobList=make([]common.Job,0)
	if getResp,err=JobMgr.kv.Get(context.TODO(),jobKey,clientv3.WithPrefix());err!=nil{
		return
	}

	for _,kvPair=range getResp.Kvs{
		if err=json.Unmarshal(kvPair.Value,&job);err==nil{
			jobList=append(jobList,job)
		}else{
			err=nil
		}
	}

	return



}
func (JobMgr *JobMgr)KillJob(name string)(err error){
	var (
		killerKey string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseid clientv3.LeaseID

	)
	killerKey=common.JOB_KILLER_DIR+name
	if leaseGrantResp,err=JobMgr.lease.Grant(context.TODO(),1);err!=nil{
		return
	}
	leaseid=leaseGrantResp.ID
	if _,err =JobMgr.kv.Put(context.TODO(),killerKey,"",clientv3.WithLease(leaseid));err!=nil{
		return
	}
	return

}
