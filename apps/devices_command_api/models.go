package main

import (
	"github.com/google/uuid"
)

type DeviceType string

const (
	DeviceTypeSensor   DeviceType = "sensor"
	DeviceTypeLighting DeviceType = "lighting"
)

type DeviceCreate struct {
	Name       string     `json:"name" binding:"required" example:"Living Room Temperature Sensor"`
	DeviceType DeviceType `json:"device_type" binding:"required" example:"sensor"`
	HomeID     uuid.UUID  `json:"home_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type DeviceAction struct {
	DeviceID uuid.UUID         `json:"device_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Action   string            `json:"action" binding:"required" example:"turn_on"`
	Params   map[string]string `json:"params,omitempty"`
}

type Response struct {
	Message string `json:"message" example:"Command sent"`
}

type CommandType string

const (
	CreateDevice CommandType = "CreateDevice"
	SendAction   CommandType = "SendAction"
)

type Command struct {
	CommandType CommandType       `json:"command_type" binding:"required" example:"CreateDevice"`
	Params      map[string]string `json:"params,omitempty"`
}
