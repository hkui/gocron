package main

import (
	"fmt"
	"time"
)

func main() {
	var (
		myChan chan map[int]int
		a map[int]int
	)
	myChan=make(chan map[int]int)
	go func() {
		for v:=range myChan{
			fmt.Println("get ",v)
			time.Sleep(3*time.Second)
		}
	}()
	for i:=0;i<10;i++{
		a= map[int]int{i:i+i}
		fmt.Println("push ",i)
		myChan<-a
		time.Sleep(1*time.Second)

	}
}
