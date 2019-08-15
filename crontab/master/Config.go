package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int64 `json:"apiReadTimeout"`
	ApiWriteTimeout int64 `json:"apiWriteTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EctdDialTimeout int64 `json:"ectdDialTimeout"`
	Webroot string `json:"webroot"`
	MongodbUri []string  `json:"mongodbUri"`
	MongodbConnectTimeout int64 `json:"mongodbConnectTimeout"`
	ShellCommand string `json:"shellCommand"`
	Yii string `json:"yii"`
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