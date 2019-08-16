package master

import (
	"strings"
)

/*
把etcd里的 php /code/yii/yii test/one  处理成 test/one

*/
func FilterCommand(commandstr string,config *Config)(string)  {
	if !config.CommandCheck{
		return commandstr
	}
	return strings.Trim(strings.Replace(commandstr,config.Yii,"",-1)," ")
}
//把 test/one 处理成完整的  php /code/yii/yii test/one
func BuildCommand(commandstr string,config *Config)(string)  {
	if !config.CommandCheck{
		return commandstr
	}
	return config.Yii+" "+strings.Trim(commandstr," ")
}
