package worker

import (
	"context"
	"crontab/common"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogSink struct {
	client *mongo.Client
	logCollection *mongo.Collection
	logChan chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}
var(
	G_logSink *LogSink
)

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)

	client,err=mongo.Connect(ctx,&options.ClientOptions{
		Hosts:G_config.MongodbUri,

	})
	if err!=nil{
		return
	}
	//检查连接
	err=client.Ping(context.TODO(),nil)
	if err!=nil{
		return
	}
	G_logSink=&LogSink{
		client:client,
		logCollection:client.Database("cron").Collection("log"),
		logChan:make(chan *common.JobLog,100),
		autoCommitChan:make(chan *common.LogBatch,1000),

	}
	go G_logSink.writeLoop()
	return

}
func (logSink *LogSink) writeLoop() {
	var (
		log *common.JobLog
		logBatch *common.LogBatch // 当前的批次
		commitTimer *time.Timer
		timeoutBatch *common.LogBatch // 超时批次
	)

	for {
		select {
		case log = <- logSink.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				// 让这个批次超时自动提交(给1秒的时间）
				commitTimer = time.AfterFunc(
					time.Duration(G_config.JobLogCommitTimeout) * time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch),
				)
			}
			fmt.Printf("追加日志:%++v\n",log)
			// 把新日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果批次满了, 就立即发送
			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				// 发送日志
				logSink.saveLogs(logBatch)
				// 清空logBatch
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <- logSink.autoCommitChan: // 过期的批次
			// 判断过期批次是否仍旧是当前的批次
			if timeoutBatch != logBatch {
				continue // 跳过已经被提交的批次
			}
			// 把批次写入到mongo中
			logSink.saveLogs(timeoutBatch)
			// 清空logBatch
			logBatch = nil
		}
	}
}

func (logSink *LogSink)Append( log *common.JobLog)  {
	select {
		case logSink.logChan<-log:  //队列满了就会阻塞,走到default,开始下一个监听

	default:

	}

}
// 批量写入日志
func (logSink *LogSink) saveLogs(batch *common.LogBatch)(err error) {
	fmt.Println("开始存储日志")
	_,err=logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
	if err!=nil {
		fmt.Println(err)
	}
	return
}
