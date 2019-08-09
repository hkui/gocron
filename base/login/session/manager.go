package session

import (
	"base/login/common"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)
var (
	G_session *Manager
)

type Manager struct {
	cookieName string //cookie的名称
	lock sync.Mutex //锁，保证并发时数据的安全一致
	provider common.Provider //管理session
	maxLifeTime int64 //超时时间
}
//在启动函数中开启GC
func init() {
	G_session, _ = NewManager("memory", "sessionid", 3600)
	go G_session.SessionGC()
}
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
//根据当前请求的cookie中判断是否存在有效的session, 不存在则创建
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session common.Session) {
	//为该方法加锁
	manager.lock.Lock()
	defer manager.lock.Unlock()
	//获取 request 请求中的 cookie 值
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name: manager.cookieName,
			Value: url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
			Path: "/",
			HttpOnly: true,
			MaxAge: int(manager.maxLifeTime)}

		http.SetCookie(w, &cookie) //将新的cookie设置到响应中
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func NewManager(providerName, cookieName string, maxLifetime int64) (*Manager, error){
	provider, ok := providers[providerName]

	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", providerName)
	}

	//返回一个 Manager 对象
	return &Manager{
		cookieName: cookieName,
		maxLifeTime: maxLifetime,
		provider: provider,
	}, nil
}
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.provider.SessionDestroy(cookie.Value)
	expiredTime := time.Now()
	newCookie := http.Cookie{
		Name: manager.cookieName,
		Path: "/", HttpOnly: true,
		Expires: expiredTime,
		MaxAge: -1,  //会话级cookie
	}
	http.SetCookie(w, &newCookie)
}
//记录该session被访问的次数
func count(w http.ResponseWriter, r *http.Request) {
	sess := G_session.SessionStart(w, r) //获取session实例
	createTime := sess.Get("createTime") //获得该session的创建时间
	if createTime == nil {
		sess.Set("createTime", time.Now().Unix())
	} else if (createTime.(int64) + 360) < (time.Now().Unix()) { //已过期
		//注销旧的session信息，并新建一个session  globalSession.SessionDestroy(w, r)
		sess = G_session.SessionStart(w, r)
	}
	count := sess.Get("countnum")
	if count == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", count.(int) + 1)
	}
}


func (manager *Manager) SessionGC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	//使用time包中的计时器功能，它会在session超时时自动调用GC方法
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() {
		manager.SessionGC()
	})
}

