package main

import (
	"fmt"
	"time"
)

func main() {
	var (
		Ch chan int
		flag chan bool
	)
	Ch=make(chan int,2)
	flag=make(chan bool)

	go func() {
		for i:=1;i<10;i++{
			Ch<-i
		}
		time.Sleep(1*time.Millisecond)
		flag<-true
	}()

	for{
		select {
			case g:=<-Ch:
			fmt.Println(g)
			case <-flag:
				goto END

		}
	}
	END:



}