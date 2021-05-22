package taillog

import (
	"fmt"
	"go_demos/my_demos/mylog/v6/kafka"

	"github.com/hpcloud/tail"
)

// 专门从日志文件收集日志模块

// TailTask  一个日志手机任务，包含单独的 tailTask  实例
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

func NewTailTask(path, topic string) (tailTask *TailTask) {
	tailTask = &TailTask{
		path:  path,
		topic: topic,
	}
	// 根据路径打开日志问阿金
	tailTask.init()
	return
}

func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开 - 文件关闭后（比如文件切割时，会关闭），重新打开
		Follow:    true,                                 // 是否跟随， (文件名发生变化后，继续读)
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 文件读取位置，
		MustExist: false,                                // 文件不存在，不报错, 会等待文件出现
		Poll:      true,                                 //
	}
	var err error
	// 打开文件
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed， err:", err)
	}

	go t.Run()
}

func (t *TailTask) Run() {
	for {
		select {
		case line := <-t.instance.Lines:
			// 发往 kafka 通道
			kafka.SendToChan(t.topic, line.Text)
		default:
		}
	}
}
