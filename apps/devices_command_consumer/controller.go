package main

import (
	"context"
	"devices_command_consumer/models"
	"devices_command_consumer/providers"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CommandController struct {
	deviceRepo        *DeviceRepository
	telemetryRepo     *TelemetryRepository
	warmHouseProvider *providers.WarmHouseProvider
}

func NewCommandController(
	deviceRepo *DeviceRepository,
	telemetryRepo *TelemetryRepository,
	warmHouseProvider *providers.WarmHouseProvider) *CommandController {
	return &CommandController{
		deviceRepo:        deviceRepo,
		telemetryRepo:     telemetryRepo,
		warmHouseProvider: warmHouseProvider,
	}
}

func (c *CommandController) ProcessCommand(ctx context.Context, message []byte) error {
	var command models.Command
	if err := json.Unmarshal(message, &command); err != nil {
		return fmt.Errorf("failed to unmarshal command: %w", err)
	}

	log.Printf("Processing command: %s", command.CommandType)

	switch command.CommandType {
	case models.CreateDevice:
		return c.handleCreateDevice(ctx, command.Params)
	case models.SendAction:
		return c.handleSendAction(ctx, command.Params)
	case models.NewTelemetry:
		return c.handleNewTelemetry(ctx, command.Params)
	default:
		return fmt.Errorf("unknown command type: %s", command.CommandType)
	}
}

func (c *CommandController) handleCreateDevice(ctx context.Context, params map[string]string) error {
	jsonData, err := json.Marshal(params)

	var device models.Device
	if err := json.Unmarshal(jsonData, &device); err != nil {
		return fmt.Errorf("failed to unmarshal device: %w", err)
	}

	provider, err := c.getProvider(ctx, &device)
	if err != nil {
		return err
	}

	outerId, err := provider.Create(&device)
	device.OuterID = outerId

	err = c.deviceRepo.Create(&device)
	if err != nil {
		return fmt.Errorf("failed to create device in database: %w", err)
	}

	log.Printf("Created device: %s (ID: %s)", device.Name, device.DeviceID)

	return nil
}

func (c *CommandController) handleSendAction(ctx context.Context, params map[string]string) error {
	deviceIDStr, exists := params["device_id"]
	if !exists {
		return fmt.Errorf("missing required parameter: device_id")
	}

	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return fmt.Errorf("invalid device_id format: %w", err)
	}

	device, err := c.deviceRepo.GetByID(deviceID)
	if err != nil {
		return fmt.Errorf("device not found: %w", err)
	}

	provider, err := c.getProvider(ctx, device)
	if err != nil {
		return err
	}

	err = provider.SendAction(device, params)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandController) handleNewTelemetry(ctx context.Context, params map[string]string) error {
	jsonData, err := json.Marshal(params)

	var telemetry models.TelemetryData
	if err := json.Unmarshal(jsonData, &telemetry); err != nil {
		return fmt.Errorf("failed to unmarshal device: %w", err)
	}

	err = c.telemetryRepo.Create(&telemetry)
	if err != nil {
		return fmt.Errorf("failed to create telemetry record: %w", err)
	}

	return nil
}

func (c *CommandController) getProvider(ctx context.Context, device *models.Device) (providers.Provider, error) {
	switch device.DeviceType {
	case models.DeviceTypeSensor:
		return c.warmHouseProvider, nil
	default:
		return nil, fmt.Errorf("unknown device type: %s", device.DeviceType)
	}
}
