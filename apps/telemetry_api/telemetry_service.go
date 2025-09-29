package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"telemetry_api/models"
	"telemetry_api/providers"
	"time"
)

type TelemetryService struct {
	deviceRepo      *DeviceRepository
	providerFactory *providers.ProviderFactory
}

func NewTelemetryService(
	deviceRepo *DeviceRepository,
	providerFactory *providers.ProviderFactory,
) *TelemetryService {
	return &TelemetryService{
		deviceRepo:      deviceRepo,
		providerFactory: providerFactory,
	}
}

func (s *TelemetryService) GetLatestTelemetryByDeviceID(deviceID uuid.UUID) (*models.TelemetryResponse, error) {
	// Get device information
	device, err := s.deviceRepo.GetByID(deviceID)
	if err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	// Get appropriate provider
	provider, err := s.providerFactory.GetProvider(device.DeviceType)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Get telemetry from provider
	telemetry, err := provider.GetTelemetry(device)
	if err != nil {
		return nil, fmt.Errorf("failed to get telemetry from provider: %w", err)
	}

	// Send to Kafka
	message := &models.Command{
		CommandType: models.NewTelemetry,
		Params: map[string]string{
			"telemetry_id": uuid.New().String(),
			"value":        strconv.FormatFloat(telemetry.Value, 'f', -1, 64),
			"device_id":    deviceID.String(),
			"created_at":   time.Now().Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal telemetry message: %w", err)
	}

	err = sendKafkaMessage("commands", messageBytes)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to send telemetry to Kafka: %v\n", err)
	}

	return &models.TelemetryResponse{
		TelemetryID: telemetry.TelemetryID,
		DeviceID:    telemetry.DeviceID,
		Value:       telemetry.Value,
		CreatedAt:   telemetry.CreatedAt,
	}, nil
}
