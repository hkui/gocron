package main

//shell命令基本执行
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	//"syscall"
)

func main() {
	var (
		cmd   *exec.Cmd
		err   error
		shell string
		out   bytes.Buffer
	)

	shell = " php /code/yii/yii"
	shell = " php /www/lepu_master/yii|grep -E '^\\s+[a-z-]+/[a-z-]+'|awk '{print $1}'"

	cmd = exec.Command("/bin/bash", "-c", shell)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(out.String())

}
