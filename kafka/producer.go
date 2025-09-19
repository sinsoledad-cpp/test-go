package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// 配置生产者
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = 3                    // 重试次数
	config.Producer.Return.Successes = true          // 返回成功的消息

	// 创建同步生产者（也可创建异步生产者）
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// 发送消息
	topic := "topic02"
	for {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("Hello Sarama! Message %s %s", time.Now(), os.Args[1])),
		}

		// 同步发送，返回分区和偏移量
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %s", err)
		} else {
			fmt.Printf("Message sent to partition %d, offset %d\n", partition, offset)
		}

		time.Sleep(time.Second)
	}

}
