package master

import (
	"net"
	"net/http"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}
var (
	G_apiServer *ApiServer
)
//保存任务接口
func handleJobSave(w http.ResponseWriter,r *http.Request)  {

}

func InitApiServer() (err error) {
	var (
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
	)
	mux=http.NewServeMux()
	mux.HandleFunc("/job/save",handleJobSave)

	//启动tcp监听

	if listener,err=net.Listen("tcp",":8070");err!=nil{
		return
	}

	//创建1个http服务
	httpServer=&http.Server{
		ReadTimeout:5*time.Second,
		WriteTimeout:5*time.Second,
		Handler:mux,
	}
	//赋值单例
	G_apiServer=&ApiServer{
		httpServer:httpServer,
	}

	go httpServer.Serve(listener)
	return



}