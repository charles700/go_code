package main

import (
	"fmt"
	"go_demos/my_demos/mylog/v5/kafka"
	"go_demos/my_demos/mylog/v5/taillog"
	"time"
)

func run() {
	// 1. 读取日志信息

	for {
		select {
		case line := <-taillog.ReadChan():
			// 2. 发送到 kafka
			kafka.SendToKafka("web_log", line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}

// LogAgent 入口程序
func main() {

	// 1 初始化 Kafka 连接
	err := kafka.Init([]string{"127.0.0.1:9092"})

	if err != nil {
		fmt.Println("init kafka faild, err:", err)
		return
	}
	fmt.Println("初始化 kafka 成功")

	// 2 打开日志文件准备收集
	err = taillog.Init("./my.log")

	if err != nil {
		fmt.Println("init taild faild, err:", err)
		return
	}

	fmt.Println("初始化 taillog 成功")

	run()
}
