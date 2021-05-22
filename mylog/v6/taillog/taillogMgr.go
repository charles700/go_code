package taillog

import (
	"fmt"
	"go_demos/my_demos/mylog/v6/etcd"
	"time"
)

// TailTaskMgr 用于管理  tailTask
type TailTaskMgr struct {
	logConf     []*etcd.LogEntry
	taskMap     map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

var tailTaskMgr *TailTaskMgr

func Init(logConf []*etcd.LogEntry) {
	tailTaskMgr = &TailTaskMgr{
		logConf:     logConf,
		taskMap:     make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry), // 无缓冲区通道
	}
	for _, logEntry := range logConf {
		NewTailTask(logEntry.Path, logEntry.Topic)
	}

	go tailTaskMgr.run()
}

// 监听自己的 newConfChan， 有了新配置到来，就处理

func (t *TailTaskMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			// 1. 配置新增
			// 2. 配置删除
			// 3. 配置变更
			fmt.Println("新的配置来了", newConf)
		default:
			time.Sleep(time.Millisecond * 50)
		}

	}
}

// 向外暴露一个 taskMgr.NewConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tailTaskMgr.newConfChan
}
