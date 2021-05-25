package kafka

import (
	"fmt"
	"go_demos/my_demos/mylog/v7/es"

	"github.com/Shopify/sarama"
)

type LogData struct {
	Data string `json:"data"`
}

func Init(addrs string, topic string) (err error) {
	consumer, err := sarama.NewConsumer([]string{addrs}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("分区列表：", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return err
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))

				// 往es 中发送消息，要保证是 json 格式
				ld := es.LogData{
					Data:  string(msg.Value),
					Topic: topic,
				}

				// 准备消息，把消息发送到 es
				es.SendToESChan(&ld)
			}
		}(pc)
		select {}
	}

	return nil
}
