package main

import (
	"fmt"
	"go_demos/my_demos/mylog/v7/conf"
	"go_demos/my_demos/mylog/v7/es"
	"go_demos/my_demos/mylog/v7/kafka"

	"gopkg.in/ini.v1"
)

func main() {

	var cfg conf.LogTransferCfg
	// 0. 加载配置文件，获取 kafka 、es 需要的链接
	err := ini.MapTo(&cfg, "./conf/cfg.ini")

	if err != nil {
		fmt.Printf("load config cfg.ini failed, err:%v \n", err)
		return
	}

	fmt.Println("0. 加载配置文件成功", cfg)

	// 1. 初始化 es
	// 1.1 提供一个往 es 写数据的函数，给 kafka 用
	err = es.Init(cfg.ESCfg.Address, cfg.ESCfg.ChanSize)
	if err != nil {
		fmt.Printf("init ES consumer failed, err:%v \n", err)
		return
	}

	fmt.Println("1. 初始化 ES 成功")
	// 2. 初始化 kafka 消费者
	// 2.1 从 topic 中的所有分区中获取新消息， 发往 es
	err = kafka.Init(cfg.KafkaCfg.Address, cfg.KafkaCfg.Topic)
	if err != nil {
		fmt.Printf("init kafka consumer failed, err:%v \n", err)
		return
	}
	fmt.Println("2. 初始化 Kafka 成功, 正在监听新消息，发往 es... ")

	select {}
}
