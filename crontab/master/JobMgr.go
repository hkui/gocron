package master

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"gocron/crontab/common"
	"math"
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
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	//建立连接
	if client,err=common.GetEtcdClient(G_config.EtcdEndpoints,G_config.EctdDialTimeout);err!=nil{
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
	//如果是更新
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
func (JobMgr *JobMgr)JobList( page int64,limit int64 )(
	jobListsRes common.JobListsRes,err error) {
	var (
		getResp *clientv3.GetResponse
		//kvPair *mvccpb.KeyValue
		jobOne common.JobListOne
		jobList []common.JobListOne
		end int64
		start int64
		sumPage int64
	)
	if page<1{
		page=1
	}

	jobList=make([]common.JobListOne,0)

	if getResp,err=JobMgr.kv.Get(
		context.TODO(),
		common.JOB_SAVE_DIR,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByModRevision,clientv3.SortDescend),
		);err!=nil{
		return
	}

	for _,kvPair:=range getResp.Kvs{
		if err=json.Unmarshal(kvPair.Value,&jobOne);err==nil{
			jobOne.ModRevision=kvPair.ModRevision
			jobList=append(jobList,jobOne)
		}else{
			err=nil
		}
	}

	start=(page-1)*limit
	if start>getResp.Count-1{
		start=0
	}
	end=start+limit
	sumPage=int64(math.Ceil(float64(getResp.Count)/float64(limit)))
	if end>getResp.Count{
		end=getResp.Count
	}

	jobListsRes=common.JobListsRes{
		Lists:jobList[start:end],
		Sum:getResp.Count,
		SumPage:sumPage,
		NowPage:page,
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
//获取一个job
func (JobMgr *JobMgr)JobOne(name string)(jobOne *common.Job,err error)  {
	var (
		jobKey string
		getResp *clientv3.GetResponse

	)
	jobKey=common.JOB_SAVE_DIR+name

	if getResp,err=JobMgr.kv.Get(context.TODO(),jobKey);err!=nil{
		return
	}
	if len(getResp.Kvs)<0{
		return
	}
	if jobOne,err=common.UnpackJob(getResp.Kvs[0].Value);err!=nil{
		return
	}
	return
}

func (jobNgr *JobMgr)CheckCronExpr(cronExpr string)(nexts []string,err error)  {
	nexts,err=common.CheckCronExpr(cronExpr)
	return
}



