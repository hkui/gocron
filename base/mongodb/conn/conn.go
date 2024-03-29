package conn

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)
var (
	session *mgo.Session
	err error
)
const url  ="139.198.170.149:27017"

func GetSession()(*mgo.Session,error){
	if session,err=mgo.Dial(url);err!=nil{
		return nil,err
	}
	return session,nil

}
func GetClient() *mongo.Client {
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
	client,err:=mongo.Connect(ctx,&options.ClientOptions{
		Hosts:[]string{url},

	})
	if err!=nil{
		log.Fatal(err.Error())
	}
	//检查连接
	err=client.Ping(context.TODO(),nil)
	if err!=nil{
		log.Fatal(err.Error())
	}
	return client
}

















