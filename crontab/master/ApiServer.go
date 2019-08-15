package master

import (
	"crontab/common"
	"encoding/json"
	"github.com/dchest/captcha"
	"github.com/gorilla/sessions"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
	store      *sessions.CookieStore
	sessid     string
}

const SESS_MAX_AGE = 1800

var (
	G_apiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir
		staticHandler http.Handler
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/captcha/", handleCaptcha)
	mux.HandleFunc("/captchaId/", handleCaptchaId)
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelte)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/one", handleJobOne)
	mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/job/workers", handleWorkerList)
	mux.HandleFunc("/job/checkexpcron", handleCheckJobCronExpr)
	mux.HandleFunc("/job/shells", handleShellList)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/logout", handleLogout)

	staticDir = http.Dir(G_config.Webroot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

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
		store:      sessions.NewCookieStore([]byte("session_cookie")),
		sessid:     "user",
	}

	go httpServer.Serve(listener)
	return

}

func checkLogin(resp http.ResponseWriter, req *http.Request) (bool, error) {
	var (
		err  error
		sess *sessions.Session
	)
	if sess, err = G_apiServer.store.Get(req, G_apiServer.sessid); err != nil {
		return false, err
	}
	sess.Options.MaxAge = SESS_MAX_AGE
	if err = sess.Save(req, resp); err != nil {
		return false, err
	}
	if len(sess.Values) < 1 || sess.Values["user"] == nil {
		return false, nil;
	}
	return true, nil

}

func handleCaptcha(resp http.ResponseWriter, req *http.Request) {
	var (
		captchaHandler http.Handler
	)
	captchaHandler = captcha.Server(160, 60)
	captchaHandler.ServeHTTP(resp, req)
}
func handleCaptchaId(resp http.ResponseWriter, req *http.Request) {
	var (
		CaptchaId string
		err       error
		bytes     []byte
		data      map[string]string
	)
	CaptchaId = captcha.NewLen(4)
	data = make(map[string]string)
	data["id"] = CaptchaId
	if bytes, err = common.BuildResponse(0, "success", data); err == nil {
		resp.Write(bytes)
	}
}

//登录验证
func handleLogin(resp http.ResponseWriter, req *http.Request) {
	var (
		err       error
		id        string
		code      string
		username  string
		password  string
		bytes     []byte
		userValid bool
		msg       string
		errno     int
		sess      *sessions.Session
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	id = req.PostForm.Get("id")     //验证码的key
	code = req.PostForm.Get("code") //输入的验证码值
	username = req.PostForm.Get("username")
	password = req.PostForm.Get("password")

	if !captcha.VerifyString(id, code) {
		errno = common.CODE_FAIL
		msg = "验证码错误"

	} else {
		userValid, _ = G_userMgr.LoginCheck(username, password)
		if userValid {
			errno = 0
			msg = "登录成功"

		} else {
			errno = common.CODE_FAIL
			msg = "用户名或密码错误"
		}

	}
	if bytes, err = common.BuildResponse(errno, msg, ""); err == nil {
		if errno == 0 {
			if sess, err = G_apiServer.store.Get(req, G_apiServer.sessid); err != nil {
				goto ERR
			}
			sess.Values["user"] = username
			sess.Options.MaxAge = SESS_MAX_AGE
			if err = sess.Save(req, resp); err != nil {
				goto ERR
			}

		}
		resp.Write(bytes)
	} else {
		goto ERR
	}
	return

ERR:
	http.Error(resp, err.Error(), http.StatusBadRequest)
}

func handleLogout(resp http.ResponseWriter, req *http.Request) {
	var (
		err error

		sess  *sessions.Session
		bytes []byte
	)
	if sess, err = G_apiServer.store.Get(req, G_apiServer.sessid); err != nil {
		goto ERR
	}
	delete(sess.Values, "user")
	sess.Options.MaxAge = SESS_MAX_AGE
	if err = sess.Save(req, resp); err != nil {
		goto ERR
	}

	bytes, _ = common.BuildResponse(common.CODE_SUCCESS, "success", "")
	resp.Write(bytes)

	return

ERR:
	http.Error(resp, err.Error(), http.StatusBadRequest)
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
		login   bool
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}
	//保存任务到etcd
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//job:{"name":"job1","command":"echo php","cronExpr":"*/5 * * * * * *"}
	postJob = req.PostForm.Get("job")

	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	if len(job.Name) < 1 || len(job.Command) < 1 || len(job.CronExpr) < 1 {

		if bytes, err = common.BuildResponse(common.CODE_PARAM_LOST, "参数缺失", nil); err == nil {
			resp.Write(bytes)
		}
		return
	}
	//保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//返回正常应答{"errno":0,"msg":"","data":{}}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", &oldJob); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(common.CODE_FAIL, err.Error(), nil); err == nil {
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
		login  bool
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	name = req.PostForm.Get("name")

	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", &oldJob); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(common.CODE_FAIL, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
}

func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		bytes []byte
		page  int64
		limit int64
		login bool

		jobListsRes common.JobListsRes
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}

	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	if pageInt, err := strconv.Atoi(req.Form.Get("page")); err == nil {
		page = int64(pageInt)
	} else {
		page = 1
	}
	if limitInt, err := strconv.Atoi(req.Form.Get("limit")); err == nil {
		limit = int64(limitInt)
	} else {
		limit = 10
	}

	if jobListsRes, err = G_jobMgr.JobList(page, limit); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", jobListsRes); err == nil {
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
	} else {
		log.Println(err)
	}
	return

ERR:
	http.Error(resp, err.Error(), http.StatusBadRequest)
}
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
		login bool
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	name = req.PostForm.Get("name")
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", nil); err == nil {
		if _, err = resp.Write(bytes); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(common.CODE_FAIL, err.Error(), nil); err == nil {
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
		job   *common.Job
		login bool
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.Form.Get("name")
	if job, err = G_jobMgr.JobOne(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", job); err == nil {
		if _, err = resp.Write(bytes); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(common.CODE_FAIL, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
}
func handleCheckJobCronExpr(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		cronExpr string
		nexts    []string
		bytes    []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	cronExpr = req.PostForm.Get("cronExpr")
	if nexts, err = G_jobMgr.CheckCronExpr(cronExpr); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", nexts); err == nil {
		if _, err = resp.Write(bytes); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(common.CODE_FAIL, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}

}

// 查询任务日志
func handleJobLog(resp http.ResponseWriter, req *http.Request) {
	var (
		err        error
		name       string // 任务名字
		skipParam  string // 从第几条开始
		limitParam string // 返回多少条
		skip       int
		limit      int
		logArr     []*common.JobLogShow
		bytes      []byte
		login      bool
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}

	// 解析GET参数
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 获取请求参数 /job/log?name=job10&skip=0&limit=10
	name = req.Form.Get("name")
	skipParam = req.Form.Get("skip")
	limitParam = req.Form.Get("limit")
	if skip, err = strconv.Atoi(skipParam); err != nil {
		skip = 0
	}
	if limit, err = strconv.Atoi(limitParam); err != nil {
		limit = 20
	}

	if logArr, err = G_logMgr.ListLog(name, skip, limit); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", logArr); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(common.CODE_FAIL, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func handleWorkerList(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		bytes   []byte
		workers []string
		login   bool
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}

	if workers, err = G_workerMgr.ListWorkers(); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", workers); err == nil {
		if _, err = resp.Write(bytes); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(common.CODE_FAIL, err.Error(), nil); err == nil {
		resp.Write(bytes)
	} else {
		log.Println(err)
	}
}

//返回可选着的shell命令
func handleShellList(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		login  bool
		bytes  []byte
		output []string
	)
	if login, err = checkLogin(resp, req); err != nil {
		goto ERR
	}

	if !login {
		bytes, _ = common.BuildResponse(common.CODE_NOT_LOGIN, common.ERR_NEED_LOGIN.Error(), "");
		if _, err = resp.Write(bytes); err != nil {
			goto ERR
		}
		return
	}

	if output,err=common.ValidShells(G_config.ShellCommand); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(common.CODE_SUCCESS, "success", output); err != nil {
		goto ERR
	}
	if _, err = resp.Write(bytes); err != nil {
		goto ERR
	}
	return
ERR:
	http.Error(resp, err.Error(), http.StatusBadRequest)

}
