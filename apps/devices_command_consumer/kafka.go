package main

import (
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaConn *kafka.Conn

func initKafkaConnection() {
	brokerAddress := "kafka:29092"

	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v", err)
		return
	}
	kafkaConn = conn

	if err := createTopicIfNotExists("commands"); err != nil {
		log.Printf("Failed to ensure topics exist: %v", err)
	}
}

func createTopicIfNotExists(topicName string) error {
	partitions, err := kafkaConn.ReadPartitions(topicName)
	if err == nil && len(partitions) > 0 {
		log.Printf("Topic %s already exists", topicName)
		return nil
	}

	topicConfig := kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = kafkaConn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %w", topicName, err)
	}

	log.Printf("Successfully created topic: %s", topicName)
	return nil
}

func createKafkaReader(brokers []string, topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
}
