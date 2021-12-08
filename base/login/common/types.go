package common

type Session interface {
	Set(key, value interface{}) error //设置Session
	Get(key interface{}) interface{}  //获取Session
	Delete(key interface{}) error     //删除Session
	SessionID() string                //当前SessionID
}
