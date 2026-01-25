package main

import (
	"fmt"
	"telemetry_api/models"

	"github.com/google/uuid"
)

type DeviceService struct {
	deviceRepo *DeviceRepository
}

func NewDeviceService(deviceRepo *DeviceRepository) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
	}
}

func (s *DeviceService) GetDevicesByHomeID(homeID uuid.UUID) ([]*models.DeviceResponse, error) {
	devices, err := s.deviceRepo.GetByHomeID(homeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	var deviceResponses []*models.DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, &models.DeviceResponse{
			DeviceID:   device.DeviceID,
			Name:       device.Name,
			DeviceType: device.DeviceType,
			HomeID:     device.HomeID,
			CreatedAt:  device.CreatedAt,
		})
	}

	return deviceResponses, nil
}
