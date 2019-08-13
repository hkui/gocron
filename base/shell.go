package main

//shell命令基本执行
import (
	"context"
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		output []byte
		err error
		shell string
	)


	shell="php /www/lepu_master/yii"
	cmd=exec.CommandContext(context.TODO(),"/bin/bash","-c",shell)

	if output,err=cmd.CombinedOutput();err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n",output)




}
