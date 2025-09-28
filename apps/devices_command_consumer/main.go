package main

import (
	"context"
	"devices_command_consumer/providers"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	deviceRepo := NewDeviceRepository(db)
	telemetryRepo := NewTelemetryRepository(db)

	warmHouseProvider := providers.NewWarmHouseProvider()

	commandController := NewCommandController(deviceRepo, telemetryRepo, warmHouseProvider)

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "commands",
		GroupID: "devices_command_consumer",
	})
	defer kafkaReader.Close()

	// Start consuming messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				message, err := kafkaReader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					continue
				}

				if err := commandController.ProcessCommand(ctx, message.Value); err != nil {
					log.Printf("Error processing command: %v", err)
				}
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
}
