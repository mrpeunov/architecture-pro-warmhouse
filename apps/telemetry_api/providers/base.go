package providers

import (
	"telemetry_api/models"
)

type Provider interface {
	GetTelemetry(device *models.Device) (*models.TelemetryData, error)
}

type BaseProvider struct{}

func (b *BaseProvider) GetTelemetry(device *models.Device) (*models.TelemetryData, error) {
	return nil, nil
}
