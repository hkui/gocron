package worker

import (
	"context"
	"crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogSink struct {
	client *mongo.Client
	logCollection *mongo.Collection
	logChan chan *common.JobLog
}
var(
	G_logSink *LogSink
)

func InitLogSink(err error)  {
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

	}

}
func (logSink *LogSink)writeLoop()  {
	var(
		log *common.JobLog
	)
	for{
		select{
		case log=<-logSink.logChan:

		}
	}
}