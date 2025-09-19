package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// GroupConsumer 自定义消费者，实现 sarama.ConsumerGroupHandler 接口
type GroupConsumer struct{}

// Setup 在分区分配完成后调用（初始化）
func (c *GroupConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 在分区被重新分配前调用（清理）
func (c *GroupConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 处理分区消息（核心逻辑）
func (c *GroupConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 遍历分区消息
	for msg := range claim.Messages() {
		// 处理消息（实际业务逻辑）
		fmt.Printf("Received message: %s (topic: %s, partition: %d, offset: %d)\n",
			string(msg.Value), msg.Topic, msg.Partition, msg.Offset)

		// 手动提交偏移量（可选，若关闭自动提交则必须手动调用）
		// 注意：需在消息处理成功后再提交，避免消息丢失
		session.MarkMessage(msg, "") // 标记消息为已处理，会在 session 结束时提交
	}
	return nil
}

func main() {
	// 配置消费组
	config := sarama.NewConfig()
	config.Version = sarama.V3_9_0_0 // 指定 Kafka 版本（需与集群版本匹配）

	// 偏移量配置：默认自动提交（可改为手动提交）
	config.Consumer.Offsets.AutoCommit.Enable = true      // 开启自动提交
	config.Consumer.Offsets.AutoCommit.Interval = 5000    // 自动提交间隔（毫秒）
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 初始偏移量（首次消费时）

	// 创建消费组
	groupID := "auto-offset-group" // 消费组唯一标识
	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{"127.0.0.1:9092"}, // Kafka 地址
		groupID,
		config,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// 要消费的主题
	topics := []string{"topic02"}
	consumer := &GroupConsumer{}

	// 启动消费循环
	go func() {
		for {
			// 持续消费（重平衡后会重新调用）
			if err := consumerGroup.Consume(context.Background(), topics, consumer); err != nil {
				log.Printf("Consume error: %v", err)
				return
			}
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down consumer group...")
}
