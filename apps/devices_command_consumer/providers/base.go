package providers

import (
	"devices_command_consumer/models"
)

type Provider interface {
	Create(device *models.Device) (string, error)
	SendAction(device *models.Device, data map[string]string) error
}

type BaseProvider struct{}

func (b *BaseProvider) Create(device models.Device) (string, error) {
	return "", nil
}

func (b *BaseProvider) SendAction(device models.Device, data map[string]string) error {
	return nil
}
