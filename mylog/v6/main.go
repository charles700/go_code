package main

import (
	"fmt"
	v6_config "go_demos/my_demos/mylog/v6/conf"
	"go_demos/my_demos/mylog/v6/etcd"
	"go_demos/my_demos/mylog/v6/kafka"
	"go_demos/my_demos/mylog/v6/taillog"
	"time"

	"gopkg.in/ini.v1"
)

var cfg = new(v6_config.AppConfig)

func run() {
	// 1. 读取日志信息

	for {
		select {
		case line := <-taillog.ReadChan():
			// 2. 发送到 kafka
			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}

// LogAgent 入口程序
func main() {

	// 0 加载配置文件
	err := ini.MapTo(cfg, "./conf/conf.ini")

	if err != nil {
		fmt.Printf("Fail to read conf file: %v", err)
		return
	}

	// 1 初始化 Kafka 连接
	err = kafka.Init([]string{cfg.KafkaConf.Addr})

	if err != nil {
		fmt.Println("init kafka faild, err:", err)
		return
	}
	fmt.Println("初始化 kafka 成功")

	// 2 初始化etcd
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Println("init etcd faild, err:", err)
		return
	}
	fmt.Println("初始化 etcd 成功")

	// 从etcd 中获取日志手机配置项的信息
	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.Key)

	if err != nil {
		fmt.Println("init etcd faild, err:", err)
		return
	}

	for index, value := range logEntryConf {
		fmt.Println("get etcd config info success, ", index, value)
	}

	// // 2 打开日志文件准备收集
	// err = taillog.Init(cfg.TaillogConf.Filename)

	// if err != nil {
	// 	fmt.Println("init taild faild, err:", err)
	// 	return
	// }

	// fmt.Println("初始化 taillog 成功")

	// run()
}
