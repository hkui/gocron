package conn

import (
	"gopkg.in/mgo.v2"
)
var (
	session *mgo.Session
	err error
)
const url  ="39.100.78.46:27017"

func GetSession()(*mgo.Session,error){
	if session,err=mgo.Dial(url);err!=nil{
		return nil,err
	}
	return session,nil

}