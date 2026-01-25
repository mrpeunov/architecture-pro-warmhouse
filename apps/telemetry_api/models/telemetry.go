package models

import (
	"time"

	"github.com/google/uuid"
)

type TelemetryData struct {
	TelemetryID uuid.UUID `json:"telemetry_id" db:"telemetry_id"`
	DeviceID    uuid.UUID `json:"device_id" db:"device_id"`
	Value       float64   `json:"value" db:"value"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type TelemetryResponse struct {
	TelemetryID uuid.UUID `json:"telemetry_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	DeviceID    uuid.UUID `json:"device_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Value       float64   `json:"value" example:"22.5"`
	CreatedAt   time.Time `json:"created_at" example:"2023-12-01T10:00:00Z"`
}

type CommandType string

const (
	NewTelemetry CommandType = "NewTelemetry"
)

type Command struct {
	CommandType CommandType       `json:"command_type" binding:"required" example:"CreateDevice"`
	Params      map[string]string `json:"params,omitempty"`
}
