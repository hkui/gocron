package main

import (
	"context"
	"fmt"
	"gocron/base/mongodb/conn"
	"gocron/base/mongodb/models"
)

func main() {
	client := conn.GetClient()

	collection := client.Database("test").Collection("stu")
	//插入1条
	/*stu:= models.Stu{
		Name:"hkui",
		Age:18,
		Sex:1,
		ClassNo:1024,
	}

	insertOneRes,err:=collection.InsertOne(context.TODO(),stu)
	if(err!=nil){
		fmt.Println(err)
	}
	fmt.Println(insertOneRes.InsertedID)*/

	//插入多条
	var manydocs []interface{}
	var oneStu models.Stu
	for i := 0; i < 10; i++ {
		oneStu = *new(models.Stu)
		oneStu.Name = fmt.Sprintf("hkui_%d", i)
		oneStu.Age = i
		oneStu.Sex = i / 2
		oneStu.ClassNo = i
		manydocs = append(manydocs, oneStu)

	}

	insertManyRes, err := collection.InsertMany(context.TODO(), manydocs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(insertManyRes.InsertedIDs)

}
