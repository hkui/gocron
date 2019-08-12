package common

import "github.com/pkg/errors"

var(
	ERR_NEED_LOGIN=errors.New("请先登录")
	ERR_LOCK_ALREADY_REQUIRED=errors.New("锁已经被占用")
	ERR_NO_LOCAL_IP_FOUND = errors.New("没有找到网卡IP")


)