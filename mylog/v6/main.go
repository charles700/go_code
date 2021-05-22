package main

import (
	"fmt"
	v6_config "go_demos/my_demos/mylog/v6/conf"
	"go_demos/my_demos/mylog/v6/etcd"
	"go_demos/my_demos/mylog/v6/kafka"
	"go_demos/my_demos/mylog/v6/taillog"
	"go_demos/my_demos/mylog/v6/utils"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

var cfg = new(v6_config.AppConfig)

var wg sync.WaitGroup

// LogAgent 入口程序
func main() {

	// 0 加载配置文件
	err := ini.MapTo(cfg, "./conf/conf.ini")

	if err != nil {
		fmt.Printf("Fail to read conf file: %v", err)
		return
	}

	// 1 初始化 Kafka 连接
	err = kafka.Init([]string{cfg.KafkaConf.Addr}, cfg.KafkaConf.MaxSize)
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

	// 2. 从etcd 中获取日志收集配置项的信息 - topic, kafka 等
	// 为了实现每个logAgent 都可以拉取自己的独有配置，所以要是实现自己的Ip 区分
	ip, err := utils.GetOutBoundIp()
	if err != nil {
		return
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ip)
	logEntryConf, err := etcd.GetConf(etcdConfKey)

	if err != nil {
		fmt.Println("init etcd faild, err:", err)
		return
	}

	// 2.初始化收集日志发往kafka
	taillog.Init(logEntryConf)

	// 3 派一个哨兵 监视 日志收集项的变化，

	newConfChan := taillog.NewConfChan()
	wg.Add(1)                                  // 从logtail中获取 chan
	go etcd.WatcConf(etcdConfKey, newConfChan) // 哨兵	发现配置有变化 及时通知 logAgent 热加载配置
	wg.Wait()
}
