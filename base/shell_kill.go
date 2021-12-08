package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	resultChan := make(chan *result, 1000)
	ctx, cancelFunc := context.WithCancel(context.TODO())
	go func() {
		cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2;ls -al")

		output, err := cmd.CombinedOutput()
		resultChan <- &result{
			err:    err,
			output: output,
		}

	}()

	time.Sleep(1 * time.Second)
	cancelFunc()

	res := <-resultChan
	fmt.Printf("err=%s\ncontent=%s\n", res.err, string(res.output))

}
