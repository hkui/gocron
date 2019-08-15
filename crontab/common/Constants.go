package common

const (
	JOB_SAVE_DIR  ="/cron/jobs/"  //保存路径目录
	JOB_KILLER_DIR  ="/cron/killer/"
	JOB_LOCK_DIR ="/cron/lock/"   //锁

	// 服务注册目录
	JOB_WORKER_DIR = "/cron/workers/"
	//任务保存
	JOB_EVENT_SAVE=1
	//任务删除
	JOB_EVENT_DELETE=2
	JOB_EVENT_KILL=3


	CODE_SUCCESS=0
	CODE_FAIL=1
	CODE_NOT_LOGIN=2
	CODE_PARAM_LOST=3
	CODE_COMMAND_INVALID=4


)
