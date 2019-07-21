package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var(
		expr *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time

	)
	if expr,err=cronexpr.Parse("*/1 * * * * * *");err!=nil{
		fmt.Println(err)
		return
	}
	//当前时间
	now=time.Now()
	//下次调度时间
	nextTime=expr.Next(now)
	//定时器差值
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("时间到 ")
	})
	time.Sleep(5*time.Second)
}

