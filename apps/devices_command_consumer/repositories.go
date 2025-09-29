package main

import (
	"database/sql"
	"devices_command_consumer/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(device *models.Device) error {
	query := `
		INSERT INTO devices (device_id, outer_id, name, device_type, home_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(
		query, device.DeviceID, device.OuterID, device.Name, device.DeviceType, device.HomeID, device.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	return nil
}

func (r *DeviceRepository) GetByID(deviceID uuid.UUID) (*models.Device, error) {
	query := `
		SELECT device_id, outer_id, name, device_type, home_id, created_at
		FROM devices
		WHERE device_id = $1`

	var device models.Device
	err := r.db.QueryRow(query, deviceID).Scan(
		&device.DeviceID,
		&device.OuterID,
		&device.Name,
		&device.DeviceType,
		&device.HomeID,
		&device.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("device not found")
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return &device, nil
}

type TelemetryRepository struct {
	db *sql.DB
}

func NewTelemetryRepository(db *sql.DB) *TelemetryRepository {
	return &TelemetryRepository{db: db}
}

func (r *TelemetryRepository) Create(telemetry *models.TelemetryData) error {
	query := `
		INSERT INTO telemetries (telemetry_id, device_id, value, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query, telemetry.TelemetryID, telemetry.DeviceID, telemetry.Value, telemetry.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create telemetry: %w", err)
	}

	return nil
}
