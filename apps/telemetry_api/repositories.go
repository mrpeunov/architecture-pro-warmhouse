package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"telemetry_api/models"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
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

func (r *DeviceRepository) GetByHomeID(homeID uuid.UUID) ([]*models.Device, error) {
	query := `
		SELECT device_id, outer_id, name, device_type, home_id, created_at
		FROM devices
		WHERE home_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, homeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var devices []*models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(
			&device.DeviceID,
			&device.OuterID,
			&device.Name,
			&device.DeviceType,
			&device.HomeID,
			&device.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, &device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return devices, nil
}
