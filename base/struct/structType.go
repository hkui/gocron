package main

import (
	"fmt"
	"reflect"
)

type MyType struct {
	Name string `json:"name"`
}
func main() {
	mt := MyType{Name: "test"}
	myType := reflect.TypeOf(mt)

	name := myType.Field(0)
	fmt.Printf("%++v\n%++v\n%++v\n",myType,myType.Field(0),name)
	tag := name.Tag.Get("json")
	println(tag)
}
