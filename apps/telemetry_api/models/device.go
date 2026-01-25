package models

import (
	"time"

	"github.com/google/uuid"
)

type DeviceType string

const (
	DeviceTypeSensor   DeviceType = "sensor"
	DeviceTypeLighting DeviceType = "lighting"
)

type Device struct {
	DeviceID   uuid.UUID  `json:"device_id" db:"device_id"`
	OuterID    string     `json:"outer_id" db:"outer_id"`
	Name       string     `json:"name" db:"name"`
	DeviceType DeviceType `json:"device_type" db:"device_type"`
	HomeID     uuid.UUID  `json:"home_id" db:"home_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

type DeviceResponse struct {
	DeviceID   uuid.UUID  `json:"device_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name       string     `json:"name" example:"Living Room Temperature Sensor"`
	DeviceType DeviceType `json:"device_type" example:"sensor"`
	HomeID     uuid.UUID  `json:"home_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedAt  time.Time  `json:"created_at" example:"2023-12-01T10:00:00Z"`
}
