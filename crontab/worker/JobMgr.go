package worker

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"gocron/crontab/common"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_jobMgr *JobMgr
)

//初始化管理器
func InitJobMgr() (err error) {
	var (
		client *clientv3.Client
	)
	//建立连接
	if client, err = common.GetEtcdClient(G_config.EtcdEndpoints, G_config.EctdDialTimeout); err != nil {
		return
	}
	G_jobMgr = &JobMgr{
		client:  client,
		kv:      client.KV,
		lease:   client.Lease,
		watcher: client.Watcher,
	}
	return
}

//监听任务变化
func (JobMgr *JobMgr) WatchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)
	if getResp, err = JobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	//初始化时载入所有的任务
	fmt.Printf("载入的任务数:%d\n", len(getResp.Kvs))
	for _, kvpair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			G_scheduler.PushJobEvent(jobEvent)
			fmt.Printf("pushJobEvent %++v\n", jobEvent.Job.Name)
		}
	}

	//从该revision向后监听变化事件
	go func() {
		watchStartRevision = getResp.Header.Revision + 1
		//监听/cron/jobs/的变化
		watchChan = JobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())

		fmt.Println("开始监听")
		for watchResp = range watchChan {
			fmt.Printf("events=%+v\n", watchResp.Events)
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					//构建event事件，
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
					G_scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE:
					jobName = common.ExtraJobName(string(watchEvent.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
					G_scheduler.PushJobEvent(jobEvent)
				default:

				}
			}
		}
	}()
	return
}
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}

//监听强杀
func (jobMgr *JobMgr) WatchKiller() {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		job        *common.Job
		jobEvent   *common.JobEvent
	)
	go func() {
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					//杀死任务只需要任务名即可
					job = &common.Job{
						Name: common.ExtraKillerName(string(watchEvent.Kv.Key)),
					}
					//构建event事件，
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_KILL, job)
					G_scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: //标记删除 不关心

				}
			}
		}

	}()
}
