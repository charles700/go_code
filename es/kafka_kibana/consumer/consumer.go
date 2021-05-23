package consumer

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// kafka 消费者实例

func InitConsumer() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)

	topic_name := "nginx_log"

	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic_name) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("partitionList", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic_name, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("kafa 消费者受到消息: Partition:%d Offset:%d Key:%v Value:%v \n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
		}(pc)

		select {}
	}

}
