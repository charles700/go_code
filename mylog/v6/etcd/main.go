package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 初始化etcd

var (
	client *clientv3.Client
)

// 需要收集的日志的配置信息
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(addr string, timeout time.Duration) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return nil
}

// 从 Etcd 中根据key 获取配置项
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		// fmt.Printf("获取到值：%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("Unmarshal etcd get key failed, err:%v\n", err)
			return
		}
	}
	return
}

func WatcConf(key string, newLogConfChan chan<- []*LogEntry) {
	rch := client.Watch(context.Background(), key) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			// 发现配置更新了 通知 tailllogMgr 中的 tailTask 的 logConf

			var newConf []*LogEntry

			if ev.Type != clientv3.EventTypeDelete {
				// 删除操作，不处理配置
				err := json.Unmarshal(ev.Kv.Value, &newConf)

				if err != nil {
					fmt.Printf("Unmarshal faild, err:%v \n", err)
					continue
				}

			}
			fmt.Println("接收到了新的配置", newConf)
			newLogConfChan <- newConf
		}
	}
}
