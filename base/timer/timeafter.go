package main

import (
	"fmt"
	"time"
)

func main() {
	var (
		i *int
		a int
	)
	a=0
	i=&a

	time.AfterFunc(
		2*time.Second,
		func(j *int ) func() {
			return func() {
				fmt.Println("time after ",*j)
			}
		}(i),

	)
	for ; *i < 10; *i++ {
		fmt.Println(*i)
		time.Sleep(1 * time.Second);
	}
}
func a() func() {
	return func(){}
}