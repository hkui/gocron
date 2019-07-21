package main

//shell命令基本执行
import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		output []byte
		err error
	)
	cmd=exec.Command("/bin/bash","-c","ls -al /")

	if output,err=cmd.CombinedOutput();err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))

}
