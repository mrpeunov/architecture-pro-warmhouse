package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

var kafkaWriter *kafka.Writer

func initKafka() {
	brokerAddress := "localhost:29092"

	writer := kafka.Writer{
		Addr:         kafka.TCP(brokerAddress),
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
	}

	kafkaWriter = &writer
}

func sendKafkaMessage(topic string, message []byte) error {
	kafkaMessage := kafka.Message{
		Topic: topic,
		Value: message,
	}

	return kafkaWriter.WriteMessages(context.Background(), kafkaMessage)
}
