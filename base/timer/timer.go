package main

import (
	"fmt"
	"time"
)

func main() {
	var (
		t *time.Timer
		//t1 *time.Ticker
	)
	t=time.NewTimer(1*time.Second)
	//t1=time.NewTicker(1*time.Second)

	go func() {
		i:=0
		for{
			select {
			case <-t.C:
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"),i)
				i++
				t.Reset(2*time.Second)


			}
		}

	}()
	time.Sleep(10*time.Second)



}
