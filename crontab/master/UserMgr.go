package master

import (
	"crontab/common"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type UserMgr struct {
	
}
var G_userMgr *UserMgr

func InitUserMgr() (err error){
	G_userMgr=&UserMgr{}
	return
}

func (userMgr *UserMgr)LoginCheck(username string,password string)( bool, error)  {
	if(username =="admin" && password=="123cron"){
		return true,nil
	}

	return false,errors.New("admin=123cron")
	//这里只做的简单的curl验证

	var (
		ret common.LoginRes
		client http.Client
		data string
		url string
		request *http.Request
		err error
		body []byte
	)
	client = http.Client{}

	data = "username="+username+"&password="+password
	url = "你的url"
	request, err = http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return false,err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8") //设置Content-Type

	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36") //设置User-Agent
	response, err := client.Do(request)

	if err != nil {
		return false,err
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)

	if err != nil {
		return false,err
	}

	if err=json.Unmarshal(body,&ret);err!=nil{
		return false,err
	}
	if ret.Status==0{
		return true,nil
	}
	return false,errors.New(ret.Note)
}

