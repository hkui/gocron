package main

//shell命令基本执行
import (
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
	)

	shell = " php /www/lepu_master/yii"

	cmd = exec.Command("/bin/bash", "-c", shell)
	/*cmd.SysProcAttr = &syscall.SysProcAttr{
		Ctty: int(os.Stdout.Fd()),
	}*/
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

}
