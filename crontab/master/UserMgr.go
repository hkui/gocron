package master

type UserMgr struct {
	
}
var G_userMgr *UserMgr

func InitUserMgr() (err error){
	G_userMgr=&UserMgr{}
	return
}

func (userMgr *UserMgr)LoginCheck(username string,password string)( bool, error)  {
	if username=="admin" && password=="123456"{
		return true,nil
	}
	return false,nil
}
func (userMgr *UserMgr)SetLoginInfo(username string) {
	
}
