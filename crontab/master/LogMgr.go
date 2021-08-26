package master

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gocron/crontab/common"
	"time"
)

// mongodb日志管理
type LogMgr struct {
	client *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	var (
		client *mongo.Client
	)

	if client,err=common.GetMongoClient(G_config.MongodbUri);err!=nil{
		return
	}

	G_logMgr = &LogMgr{
		client: client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

// 查看任务日志
func (logMgr *LogMgr) ListLog(name string, skip int, limit int) (logArr []*common.JobLogShow, err error){
	var (
		filter *common.JobLogFilter
		logSort *common.SortLogByStartTime
		cursor *mongo.Cursor
		jobLog *common.JobLog
		findoptions options.FindOptions
	)

	// len(logArr)
	logArr = make([]*common.JobLogShow, 0)

	// 过滤条件
	filter = &common.JobLogFilter{JobName: name}

	// 按照任务开始时间倒排
	logSort = &common.SortLogByStartTime{SortOrder: -1}
	skip64:=int64(skip)
	limit64:=int64(limit)
	findoptions=options.FindOptions{
		Sort:logSort,
		Skip:&skip64,
		Limit:&limit64,
	}
	// 查询
	if cursor, err = logMgr.logCollection.Find(context.TODO(),filter,&findoptions); err != nil {
		return
	}
	// 延迟释放游标
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		// 反序列化BSON
		if err = cursor.Decode(jobLog); err != nil {
			continue // 有日志不合法
		}

		logArr = append(logArr, &common.JobLogShow{
			JobName:jobLog.JobName,
			Command:jobLog.Command,
			Err:jobLog.Err,
			Output: jobLog.Output,
			PlanTime:common.TimeToStr(time.Unix(0, jobLog.PlanTime*int64(time.Millisecond))),
			ScheduleTime:common.TimeToStr(time.Unix(0, jobLog.ScheduleTime*int64(time.Millisecond))),
			StartTime:common.TimeToStr(time.Unix(0, jobLog.StartTime*int64(time.Millisecond))),
			EndTime:common.TimeToStr(time.Unix(0, jobLog.EndTime*int64(time.Millisecond))),

		})
	}
	return
}


