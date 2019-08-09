package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
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
func login(resp http.ResponseWriter,  r *http.Request)  {
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	cookie := http.Cookie{Name: "username", Value: "hkui", Expires: expiration}
	http.SetCookie(resp, &cookie)



	resp.Write(([]byte)("123"));
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
