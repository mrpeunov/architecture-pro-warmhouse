package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var kafkaWriter *kafka.Writer
var kafkaConn *kafka.Conn

func initKafka() {
	brokerAddress := "kafka:9092"

	writer := kafka.Writer{
		Addr:         kafka.TCP(brokerAddress),
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
	}

	kafkaWriter = &writer

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

func sendKafkaMessage(topic string, message []byte) error {
	kafkaMessage := kafka.Message{
		Topic: topic,
		Value: message,
	}

	return kafkaWriter.WriteMessages(context.Background(), kafkaMessage)
}
