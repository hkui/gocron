package main

import (
	"base/mongodb/conn"
	"fmt"
	"gopkg.in/mgo.v2"
)
type stu struct {
	Name string `bson:"name"`
	Age int `bson:"age"`
	Sex int
	ClassNo int `bson:"class_no"`
}


func main() {
	var (
		session *mgo.Session
		err error
	)
	session,err=conn.GetSession()
	defer session.Close()
	if err!=nil{
		fmt.Println(err)
		return
	}
	c:=session.DB("test").C("stu")
	//单条插入
	/*for i:=0;i<10;i++{
		stu1:=new (stu)
		stu1.Name="abc"
		stu1.Age=i
		stu1.Sex=1
		stu1.ClassNo=2
		if err=c.Insert(stu1);err!=nil{
			fmt.Println(err)
		}

	}*/
	//多条插入
	var manyStus [] *stu
	for i:=0;i<10;i++{
		stu1:=new (stu)
		stu1.Name="abc"
		stu1.Age=i
		stu1.Sex=1
		stu1.ClassNo=2
		manyStus=append(manyStus,stu1)
	}

	resErr:=c.Insert(manyStus)
	if resErr!=nil{
		fmt.Println(resErr)
	}

	//查询
	//var res []interface{}

	//c.Find(bson.M{"name":"abc"}).All(&res) //根据name查询
	//c.Find(bson.M{"age":bson.M{"$gte":5}}).All(&res) //根据name查询
	//for _,v:=range res{
	//	fmt.Println(v)
	//}
	//根据id查询
/*
	idStr := "5d3ac4356ac9e2574daf3cdc"
	objectId := bson.ObjectIdHex(idStr)
	var one stu
	c.Find(bson.M{"_id":objectId}).One(&one)
	fmt.Printf("%+v",one)*/















}
