package main

import (
	"fmt"
	"time"
)

func main() {

	t, err := msToTime(1564886970004)
	if err != nil {
		return
	}
	fmt.Println(t, "-", int64(time.Millisecond))

}
func msToTime(msInt int64) (time.Time, error) {

	tm := time.Unix(0, msInt*int64(time.Millisecond))

	//fmt.Println(tm.Format("2006-02-01 15:04:05.000"))

	return tm, nil
}
