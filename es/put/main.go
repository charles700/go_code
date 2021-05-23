package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		panic(err)
	}

	p1 := Person{Name: "zhangsan", Age: 22, Married: false}

	//  插入一条记录
	put1, err := client.Index().
		Index("user"). // 指定数据库
		BodyJson(p1).  //
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	// 输出插入的数据的id 和 索引库
	fmt.Printf("INSERT -- Indexs user %s to index %s , type %s \n", put1.Id, put1.Index, put1.Type)
}
