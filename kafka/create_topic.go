package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

// 检查主题是否存在（旧版API）
func topicExists(client sarama.Client, topic string) (bool, error) {
	// 获取所有主题
	topics, err := client.Topics()
	if err != nil {
		return false, fmt.Errorf("获取主题列表失败: %v", err)
	}

	// 检查目标主题是否在列表中
	for _, t := range topics {
		if t == topic {
			return true, nil
		}
	}
	return false, nil
}

// 创建主题（旧版API）
func createTopicIfNotExists(brokers []string, topic string, partitions int32, replicationFactor int16) error {
	// 配置客户端
	config := sarama.NewConfig()
	config.Net.DialTimeout = 5 * time.Second

	// 创建基础客户端（旧版没有专门的AdminClient）
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return fmt.Errorf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 检查主题是否存在
	exists, err := topicExists(client, topic)
	if err != nil {
		return fmt.Errorf("检查主题存在性失败: %v", err)
	}

	if exists {
		log.Printf("主题 %s 已存在，无需创建", topic)
		return nil
	}

	// 构建创建主题的请求
	request := sarama.CreateTopicsRequest{
		TopicDetails: map[string]*sarama.TopicDetail{
			topic: {
				NumPartitions:     partitions,
				ReplicationFactor: replicationFactor,
				ConfigEntries:     make(map[string]*string),
			},
		},
		Timeout: 10 * time.Second, // 10秒超时
	}

	// 通过客户端发送创建请求（旧版需手动选择broker）
	broker, err := client.Controller()
	if err != nil {
		return fmt.Errorf("无法获取控制器broker")
	}

	// 发送请求
	response, err := broker.CreateTopics(&request)
	if err != nil {
		return fmt.Errorf("发送创建请求失败: %v", err)
	}

	// 检查响应结果
	for topic, err := range response.TopicErrors {
		fmt.Println(topic, err)
	}

	log.Printf("主题 %s 创建成功，分区数: %d，副本数: %d", topic, partitions, replicationFactor)
	return nil
}

func main() {
	brokers := []string{"127.0.0.1:9092"}
	topicName := "topic02"
	partitionCount := int32(3)
	replicationFactor := int16(1)

	err := createTopicIfNotExists(brokers, topicName, partitionCount, replicationFactor)
	if err != nil {
		log.Fatalf("操作失败: %v", err)
	}
	log.Println("操作完成")
}
