package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CommandType string

const (
	CreateDevice CommandType = "CreateDevice"
	SendAction   CommandType = "SendAction"
	NewTelemetry CommandType = "NewTelemetry"
)

type Command struct {
	CommandType CommandType       `json:"command_type" binding:"required" example:"CreateDevice"`
	Params      map[string]string `json:"params,omitempty"`
}

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

type TelemetryData struct {
	TelemetryID uuid.UUID   `json:"telemetry_id" db:"telemetry_id"`
	DeviceID    uuid.UUID   `json:"device_id" db:"device_id"`
	Value       json.Number `json:"value" db:"value"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

type TelemetryCreate struct {
	DeviceID uuid.UUID `json:"device_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Value    float64   `json:"value" binding:"required" example:"22.5"`
}
