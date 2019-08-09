package main

import (
	"base/login/session"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)




func main() {
	session.Init()
	var(
		mux *http.ServeMux
		listener net.Listener
		err error
		httpServer *http.Server
	)

	mux=http.NewServeMux()
	mux.HandleFunc("/login",login)
	mux.HandleFunc("/get",get)

	if listener,err=net.Listen("tcp",":8888");err!=nil{
		panic(err)
	}
	httpServer=&http.Server{
		Handler:mux,
	}
	err=httpServer.Serve(listener)
	if err!=nil{
		panic(err)
	}

}
func login(w http.ResponseWriter,  r *http.Request)  {
	sess := G_session.SessionStart(w, r)
	r.ParseForm()
	name := sess.Get("username")
	if name != nil {
		sess.Set("username", r.Form["username"]) //将表单提交的username值设置到session中
	}
	w.Write(([]byte)("123"));
}

func get(resp http.ResponseWriter,  r *http.Request)()  {
	var (
		cookie *http.Cookie
		err error
		key string
	)
	if err:=r.ParseForm();err!=nil{
		fmt.Println(err)
		return
	}
	key=r.Form.Get("key")

	if cookie,err=r.Cookie(key);err!=nil{
		fmt.Println(err)
		return
	}
	resp.Write([]byte(cookie.Value))

}



