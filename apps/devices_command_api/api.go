package main

import (
	"encoding/json"
	"maps"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeviceHandler handles HTTP requests for device operations
type DeviceHandler struct{}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

// RegisterRoutes registers all device routes
func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
	devices := router.Group("/devices")
	{
		devices.POST("", h.CreateDevice)
		devices.POST("/:device_id/action", h.SendAction)
	}
}

// CreateDevice creates a new device
// @Summary Create a new device
// @Description Create a new smart home device
// @Tags devices
// @Accept json
// @Produce json
// @Param device body DeviceCreate true "Device creation data"
// @Success 201 {object} Command
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var deviceCreate DeviceCreate
	if err := c.ShouldBindJSON(&deviceCreate); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	command := Command{
		CommandType: CreateDevice,
		Params: map[string]string{
			"device_id":   uuid.New().String(),
			"device_type": string(deviceCreate.DeviceType),
			"name":        deviceCreate.Name,
			"home_id":     deviceCreate.HomeID.String(),
			"created_at":  time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	commandData, err := json.Marshal(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: "Failed to marshal command"})
		return
	}

	err = sendKafkaMessage("commands", commandData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: "Can't send command"})
		return
	}

	c.JSON(http.StatusOK, command)
}

// SendAction sends a action to a device
// @Summary Send command to device
// @Description Send a command to a specific device
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param command body DeviceAction true "Device command"
// @Success 200 {object} Command
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /devices/{device_id}/action [post]
func (h *DeviceHandler) SendAction(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	_, err := uuid.Parse(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "invalid device ID"})
		return
	}

	var action DeviceAction
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	params := map[string]string{
		"device_id":  deviceIDStr,
		"action":     action.Action,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	maps.Insert(params, maps.All(action.Params))

	command := Command{
		CommandType: SendAction,
		Params:      params,
	}

	commandData, err := json.Marshal(command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: "Can't send command"})
		return
	}

	err = sendKafkaMessage("commands", commandData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Message: "Can't send command"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "Command sent"})
}
