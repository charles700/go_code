package es

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

type LogData struct {
	Data  string
	Topic string
}

var client *elastic.Client
var ch chan *LogData

func Init(address string, chanSize int) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		fmt.Printf("init es client failed, err:%v\n", err)
		return
	}

	ch = make(chan *LogData, chanSize)
	go SendToES()

	return
}

func SendToES() error {
	for {
		select {
		case msg := <-ch:
			//  插入一条记录
			put1, err := client.Index().
				Index(msg.Topic). // 指定数据库
				Type("xxx").      // 指定表, 可以传 不同的id, 按照机器来划分 log
				BodyJson(msg).    // 发送的数据
				Do(context.Background())

			if err != nil {
				fmt.Printf("send msg to es failed, err:%v\n", err)
				continue
			}
			fmt.Printf("INSERT ES SUCCESS -- Indexs %s %s to index %s , type %s \n", msg.Topic, put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}

func SendToESChan(data *LogData) {
	ch <- data
}
