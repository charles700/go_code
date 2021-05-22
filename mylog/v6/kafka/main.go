package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// 把日志 写入 kafka

type logData struct {
	topic string
	data  string
}

// 声明全局连接 ,kafka 的 生产者客户端
var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return err
	}

	logDataChan = make(chan *logData, maxSize)

	// 开启后台，从通道中取数据，发往kafka
	go SendToKafka()

	return nil
}

func SendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)

			// 发送到 Kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("SendToKafka success pid:%v offset:%v topic:%v \n", pid, offset, ld.topic)

		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}

func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}
