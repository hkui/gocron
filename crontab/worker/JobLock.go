package worker

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"gocron/crontab/common"
)

//分布式锁,txn事务
type JobLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string
	cancelFunc context.CancelFunc //用于取消自动续租
	leaseId    clientv3.LeaseID
	isLocked   bool //是否上锁成功
}

//初始化一把锁

func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		jobName: jobName,
		kv:      kv,
		lease:   lease,
	}
	return
}
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		canceCtx       context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)
	//1 创建租约
	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	//context 用于取消自动续租
	canceCtx, cancelFunc = context.WithCancel(context.TODO())
	leaseId = leaseGrantResp.ID

	//2自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(canceCtx, leaseId); err != nil {
		goto FAIL

	}
	//处理续租应答协程
	go func() {
		var (
			keepAliveResp *clientv3.LeaseKeepAliveResponse
		)
		select {
		case keepAliveResp = <-keepRespChan:
			if keepAliveResp == nil {
				goto END
			}
		}
	END:
	}()

	//4 创建事务
	txn = jobLock.kv.Txn(context.TODO())
	lockKey = common.JOB_LOCK_DIR + jobLock.jobName
	//5事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))
	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}
	//抢锁成功
	if !txnResp.Succeeded { //锁被占用
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto FAIL
	}
	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc
	jobLock.isLocked = true
	return

FAIL:
	cancelFunc()
	jobLock.lease.Revoke(context.TODO(), leaseId) //释放租约

	return

}
func (jobLock *JobLock) Unlock() {
	if jobLock.isLocked {
		jobLock.cancelFunc() //取消自动续租协程
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseId)
	}

}
