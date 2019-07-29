package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var(
		mycronExp string
		exp *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time
	)
	mycronExp="2 * * * *";
	exp,err=cronexpr.Parse(mycronExp)
	if err!=nil{
		fmt.Println(err)
		return
	}

	now=time.Now()
	nextTime=exp.Next(now)

	for i:=0;i<20;i++{
		if i==0{
			nextTime=exp.Next(now)
		}else{
			nextTime=exp.Next(nextTime)
		}
		fmt.Println(nextTime.Format("2006-01-02 15:04:05"))
	}
}
