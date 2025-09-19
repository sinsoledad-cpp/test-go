package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

func main() {
	// 配置消费者
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 从最早的消息开始消费

	// 2. 创建基础消费者（非消费组模式）
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	// 3. 指定要消费的主题和分区
	topic := "test-topic"
	partition := int32(0) // 分区编号（从0开始）

	// 4. 获取分区消费者（从指定偏移量开始消费）
	// 第三个参数为起始偏移量：可以是 sarama.OffsetOldest / sarama.OffsetNewest 或具体数值
	pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %s", err)
	}
	defer pc.Close()

	for msg := range pc.Messages() {
		// 成功接收到消息
		fmt.Printf("%s %s %d %d\n", msg.Topic, msg.Value, msg.Partition, msg.Offset)
	}
}
