package providers

import (
	"fmt"
	"telemetry_api/models"
)

type ProviderFactory struct {
	warmHouseProvider *WarmHouseProvider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		warmHouseProvider: NewWarmHouseProvider(),
	}
}

func (f *ProviderFactory) GetProvider(deviceType models.DeviceType) (Provider, error) {
	switch deviceType {
	case models.DeviceTypeSensor:
		return f.warmHouseProvider, nil
	default:
		return nil, fmt.Errorf("unknown device type: %s", deviceType)
	}
}
