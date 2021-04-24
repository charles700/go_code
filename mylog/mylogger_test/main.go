package main

import (
	mylogger "go_demos/my_demos/mylog/v1"
)

// 测试自己写的日志库
func main() {
	log := mylogger.NewLog("debug")

	log.Debug("debug= %s", "错误日志")
	log.Trace("Trace")
	log.Info("info")
	log.Warning("Warning")
	log.Error("Error")
	log.Fatal("Fatal")
}
