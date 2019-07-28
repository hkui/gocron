package worker

import (
	"context"
	"crontab/common"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	watcher clientv3.Watcher
}
var (
	G_jobMgr *JobMgr
)
//初始化管理器
func InitJobMgr() (err error) {
	var(
		config clientv3.Config
		client *clientv3.Client
	)
	config=clientv3.Config{
		Endpoints:G_config.EtcdEndpoints,//集群地址
		DialTimeout:time.Millisecond*time.Duration(G_config.EctdDialTimeout),
	}
	//建立连接

	if client,err=clientv3.New(config);err!=nil{
		return
	}
	G_jobMgr=&JobMgr{
		client:client,
		kv:client.KV,
		lease:client.Lease,
		watcher:client.Watcher,
	}
	return

}

//监听任务变化
func (JobMgr *JobMgr)WatchJobs()(err error)  {
	var(
		getResp *clientv3.GetResponse
		kvpair *mvccpb.KeyValue
		job *common.Job
		watchStartRevision int64
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobName string
		jobEvent *common.JobEvent

	)
	if getResp,err=JobMgr.kv.Get(context.TODO(),common.JOB_SAVE_DIR,clientv3.WithPrefix());err!=nil{
		return
	}
	for _,kvpair=range getResp.Kvs{
		if job,err=common.UnpackJob(kvpair.Value);err==nil{
			jobEvent=common.BuildJobEvent(common.JOB_EVENT_SAVE,job)
			G_scheduler.PushJobEvent(jobEvent)
			fmt.Printf("%++v\n",jobEvent)

		}
	}
	//从该revision向后监听变化事件
	go func() {
		watchStartRevision=getResp.Header.Revision+1
		//监听/cron/jobs/的冰花
		watchChan=JobMgr.watcher.Watch(context.TODO(),common.JOB_SAVE_DIR,clientv3.WithPrefix())

		for watchResp=range watchChan{
			for _,watchEvent=range watchResp.Events{
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job,err=common.UnpackJob(watchEvent.Kv.Value);err!=nil{
						continue
					}
					//构建event事件，
					jobEvent=common.BuildJobEvent(common.JOB_EVENT_SAVE,job)
					G_scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE:
					jobName=common.ExtraJobName(string(watchEvent.Kv.Key))
					job=&common.Job{
						Name:jobName,
					}
					jobEvent=common.BuildJobEvent(common.JOB_EVENT_DELETE,job)
					G_scheduler.PushJobEvent(jobEvent)
				default:

				}
			}
		}
	}()
	return
}

