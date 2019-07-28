package main

import (
	"fmt"
)
type worker struct {
	in chan int
	done chan  bool
}
func createWorker(id int)worker{
	w:=worker{
		in:make(chan int),
		done:make(chan  bool),
	}
	go dowork(id,w.in,w.done)
	return w
}

func dowork(id int,in chan int,done chan  bool )  {
	for n:=range in {
		fmt.Printf("worker %d reveived %c\n",id,n)
		go func() {
			done<-true
		}()
	}
}
const GNUM=3
func chanDemo()  {
	var workers [GNUM] worker
	for i:=0;i<GNUM;i++{
		workers[i]=createWorker(i)
	}

	for i,worker:=range workers{
		fmt.Println("push a",i)
		worker.in<-'a'+i
	}
	for i,worker:=range workers{
		fmt.Println("push A",i)
		worker.in<-'A'+i
	}
	for _,worker:=range  workers{
		<-worker.done
		<-worker.done
	}

}

func main() {
	chanDemo()

}