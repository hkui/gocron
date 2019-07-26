package models

type Stu struct {
	Name string `bson:"name"`
	Age int `bson:"age"`
	Sex int
	ClassNo int `bson:"class_no"`
}
