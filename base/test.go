package main

import "fmt"

func main() {
	var (
		stringArr []string
	)
	stringArr = []string{"a", "b"}
	for k, v := range stringArr {
		fmt.Println(k, v)
	}

}
