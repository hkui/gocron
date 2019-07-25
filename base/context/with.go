package main

import (
	"context"
	"fmt"
	"time"
)

func inc(a int)int  {
	res:=a+1
	time.Sleep(1*time.Second)
	return res
}
func Add(ctx context.Context,a int)int  {
	res:=0
	for i:=0;i<a ;i++  {
		res=inc(res)

		select {
			case <-ctx.Done():     //Done会返回一个channel，当该context被取消的时候，该channel会被关闭，同时对应的使用该context的routine也应该结束并返回
				fmt.Printf("get channel %d\n",i)
				return -1
			default:
				fmt.Println(i)
		}

	}

	return res

}
func main() {
	a := 1
	b := 2
	timeout := 2 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	res := Add(ctx, 3)

	fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
}
