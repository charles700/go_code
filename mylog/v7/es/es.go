package es

import (
	"context"
	"fmt"
	"strings"

	"github.com/olivere/elastic"
)

var client *elastic.Client

func Init(address string) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		fmt.Printf("init es client failed, err:%v\n", err)
		return
	}
	return
}

func SendToES(index string, data interface{}) error {

	//  插入一条记录
	put1, err := client.Index().
		Index(index).   // 指定数据库
		Type("xxx").    // 指定表, 可以传 不同的id, 按照机器来划分 log
		BodyJson(data). //
		Do(context.Background())

	if err != nil {
		fmt.Printf("send msg to es failed, err:%v\n", err)
		return err
	}

	// 输出插入的数据的id 和 索引库
	fmt.Printf("INSERT -- Indexs %s %s to index %s , type %s \n", index, put1.Id, put1.Index, put1.Type)

	return nil
}
