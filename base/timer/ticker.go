package main

import (
	"fmt"
	"time"
)

func main() {
	var (


		push chan time.Time

	)
	push=make(chan time.Time)

	go func() {
		for{
			time.Sleep(3*time.Second)
			fmt.Println("get:",<-push)
		}
	}()


	for{
		select {
		case push<-time.Now():
		default:

			fmt.Println("阻塞了",time.Now())
			time.Sleep(1*time.Second)


	}
	}




}
