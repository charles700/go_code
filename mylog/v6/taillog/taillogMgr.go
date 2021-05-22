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
		tailTask := NewTailTask(logEntry.Path, logEntry.Topic)
		// 存储一份map用来做新增、更新、删除
		key := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tailTaskMgr.taskMap[key] = tailTask
	}

	go tailTaskMgr.run()
}

// 监听自己的 newConfChan， 有了新配置到来，就处理

func (t *TailTaskMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			// 配置新增、变更
			for _, config := range newConf {
				key := fmt.Sprintf("%s_%s", config.Path, config.Topic)
				_, ok := t.taskMap[key]
				if ok {
					// 原来就有
					continue
				} else {
					// 新增task
					tailTask := NewTailTask(config.Path, config.Topic)
					t.taskMap[key] = tailTask
				}
			}
			// 处理配置删除
			for _, c1 := range t.logConf {
				isDel := true
				for _, c2 := range newConf {
					if c1.Path == c2.Path && c1.Topic == c2.Topic {
						isDel = false
						continue
					}
				}

				if isDel {
					// 把c1 对应的 task 停止掉
					key := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					_, ok := t.taskMap[key]
					if ok {
						t.taskMap[key].cancelFunc()
					}
				}
			}
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
