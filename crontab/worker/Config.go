package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {

	EtcdEndpoints []string `json:"etcdEndpoints"`
	EctdDialTimeout int64 `json:"ectdDialTimeout"`
	MongodbUri []string  `json:"mongodbUri"`
	MongodbConnectTimeout int64 `json:"mongodbConnectTimeout"`
	JobLogCommitTimeout int64 `json:"jobLogCommitTimeout"`
	JobLogBatchSize int `json:"jobLogBatchSize"`

}
var (
	G_config *Config
)

func InitConfig(filename string)(err error)  {
	var (
		content []byte
		conf Config
	)
	if content,err=ioutil.ReadFile(filename);err!=nil{
		return
	}

	if err=json.Unmarshal(content,&conf);err!=nil{
		return
	}
	G_config=&conf
	return

}