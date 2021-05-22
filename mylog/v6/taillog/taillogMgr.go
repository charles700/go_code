package taillog

import "go_demos/my_demos/mylog/v6/etcd"

// TailTaskMgr 用于管理  tailTask
type TailTaskMgr struct {
	logConf []*etcd.LogEntry
	// taskMap map[string]*TailTask
}

var tailTaskMgr *TailTaskMgr

func Init(logConf []*etcd.LogEntry) {
	tailTaskMgr = &TailTaskMgr{
		logConf: logConf,
	}
	for _, logEntry := range logConf {
		NewTailTask(logEntry.Path, logEntry.Topic)
	}
}
