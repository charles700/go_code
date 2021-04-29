package main

import (
	mylogger "go_demos/my_demos/mylog/v3"
	"time"
)

// 定义全局logger
var log mylogger.Logger

// 测试自己写的日志库
func main() {
	// 实现接口后，复制 文件日志 和 终端日志
	// log = mylogger.NewConsoleLogger("debug")
	log = mylogger.NewFileLogger("debug", "./", "debug.log", 10*1024)

	for {
		// log.Debug("debug= %s", "错误日志")
		// log.Trace("Trace")
		log.Info("info")
		// log.Warning("Warning")
		// log.Error("Error")
		// log.Fatal("Fatal")
		time.Sleep(time.Second)
	}
}
