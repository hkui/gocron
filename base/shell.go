package main

//shell命令基本执行
import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	var (
		cmd *exec.Cmd
		output []byte
		err error
		stringout string
		stringArr []string
	)
	//cmd=exec.Command("/bin/bash","-c","ls -al /")
	cmd=exec.CommandContext(context.TODO(),"/bin/bash","-c","php /code/yii/yii|grep -E '[a-z-]+/[a-z-]+'|awk '{print $1}'")

	if output,err=cmd.CombinedOutput();err!=nil{
		fmt.Println(err)
		return
	}
	stringout=string(output)

	//fmt.Println(stringout)
	stringArr=strings.Split(stringout,"\n")
	fmt.Println(stringArr,len(stringArr))


}
