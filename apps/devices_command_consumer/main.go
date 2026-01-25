package main

import (
	"context"
	"devices_command_consumer/providers"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Kafka connection and ensure topics exist
	initKafkaConnection()
	if kafkaConn != nil {
		defer kafkaConn.Close()
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	deviceRepo := NewDeviceRepository(db)
	telemetryRepo := NewTelemetryRepository(db)

	warmHouseProvider := providers.NewWarmHouseProvider()

	commandController := NewCommandController(deviceRepo, telemetryRepo, warmHouseProvider)

	kafkaReader := createKafkaReader([]string{"kafka:9092"}, "commands", "devices_command_consumer")
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
