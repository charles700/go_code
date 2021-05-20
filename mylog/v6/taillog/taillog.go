package taillog

import (
	"fmt"

	"github.com/hpcloud/tail"
)

// 专门从日志文件收集日志模块

var (
	tailObj *tail.Tail
)

func ReadChan() <-chan *tail.Line {
	return tailObj.Lines
	// var (
	// 	msg *tail.Line // 行数据
	// 	ok  bool
	// )

	// for {
	// 	msg, ok = <-tailObj.Lines // 一行行读取日志
	// 	if !ok {
	// 		fmt.Printf("tail file close reopen, filename: %s \n", tailObj.Filename)
	// 		time.Sleep(time.Second)
	// 		continue
	// 	}
	// 	fmt.Println("msg:", msg.Text)
	// 	logChan <- msg.Text
	// }
}

func Init(fileName string) (err error) {

	config := tail.Config{
		ReOpen:    true,                                 // 重新打开 - 文件关闭后（比如文件切割时，会关闭），重新打开
		Follow:    true,                                 // 是否跟随， (文件名发生变化后，继续读)
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 文件读取位置，
		MustExist: false,                                // 文件不存在，不报错, 会等待文件出现
		Poll:      true,                                 //
	}

	// 打开文件
	tailObj, err = tail.TailFile(fileName, config)

	if err != nil {
		fmt.Println("tail file failed， err:", err)
		return err
	}

	return nil
}
