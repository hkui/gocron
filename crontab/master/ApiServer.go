package master

import (
	"crontab/common"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}
var (
	G_apiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
		staticDir http.Dir
		staticHandler http.Handler
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelte)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/one", handleJobOne)
	mux.HandleFunc("/job/checkexpcron", handleCheckJobCronExpr)

	staticDir=http.Dir(G_config.Webroot)
	staticHandler=http.FileServer(staticDir)
	mux.Handle("/",http.StripPrefix("/",staticHandler))

	//启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	//创建1个http服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	//赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)
	return

}

//保存任务接口
//POST job={"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	//保存任务到etcd
	//1.解析post表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//job:{"name":"job1","command":"echo php","cronExpr":"*/5 * * * * * *"}
	postJob = req.PostForm.Get("job")

	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	if len(job.Name)<1 ||len(job.Command)<1||len(job.CronExpr)<1{
		if bytes, err = common.BuildResponse(-1, "参数缺失", nil); err == nil {
			resp.Write(bytes)
		}
		return
	}
	//保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//返回正常应答{"errno":0,"msg":"","data":{}}
	if bytes, err = common.BuildResponse(0, "success", &oldJob); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

//post /job/delete name=job1
func handleJobDelte(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	name = req.PostForm.Get("name")

	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", &oldJob); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
}

func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []common.Job
		err     error
		bytes   []byte
	)
	if jobList, err = G_jobMgr.JobList(); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {

		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}

	} else {
		log.Println(err)
	}

ERR:
	log.Println(err)

}
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	name = req.PostForm.Get("name")
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		if _,err= resp.Write(bytes);err!=nil{
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
}
func handleJobOne(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
		job *common.Job
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	name = req.Form.Get("name")
	if job,err = G_jobMgr.JobOne(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", job); err == nil {
		if _,err= resp.Write(bytes);err!=nil{
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
}
func handleCheckJobCronExpr(resp http.ResponseWriter, req *http.Request)  {
	var (
		err error
		cronExpr  string
		nexts []string
		bytes []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	cronExpr = req.PostForm.Get("cronExpr")
	if nexts,err = G_jobMgr.CheckCronExpr(cronExpr); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", nexts); err == nil {
		if _,err= resp.Write(bytes);err!=nil{
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		 log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println(err)
	}

}

