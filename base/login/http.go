package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net"
	"net/http"

)


var store = sessions.NewFilesystemStore("./",[]byte("session_file"))
//var store=sessions.NewCookieStore([]byte("session_cookie"))

const sess_name  = "user"


func main() {
	var(
		mux *http.ServeMux
		listener net.Listener
		err error
		httpServer *http.Server
	)


	mux=http.NewServeMux()
	mux.HandleFunc("/set",set)
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

func set(w http.ResponseWriter,  r *http.Request)  {
	var(
		err error
		key string
		value string
		ses *sessions.Session
	)
	err=r.ParseForm()
	if(err!=nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key=r.Form.Get("key")
	value=r.Form.Get("value")
	if len(key)<1 || len(value)<1{
		http.Error(w,fmt.Errorf("参数缺失","参数缺失").Error(),http.StatusBadRequest)
		return
	}
	if ses,err=store.Get(r,sess_name);err!=nil{
		fmt.Println(err)
		return
	}
	ses.Values[key]=value
	ses.Options.MaxAge=20

	// Save it before we write to the response/return from the handler.
	err=ses.Save(r, w)
	if err!=nil{
		fmt.Println(err)
	}else{
		fmt.Printf("set %s=%s \n",key,value)
	}
}

func get(w http.ResponseWriter,  r *http.Request)()  {
	var(
		err error

		out string
		ses *sessions.Session
	)
	err=r.ParseForm()
	if(err!=nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ses,err=store.Get(r,sess_name);err!=nil{
		fmt.Println(err)
		return
	}


	for k,v:=range ses.Values{
		out+=fmt.Sprintf("%s=%s\n",k,v)
	}
	fmt.Println(out)
	w.Write([]byte(out))



}



