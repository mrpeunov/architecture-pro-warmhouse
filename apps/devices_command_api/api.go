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
	devices.Use(AuthMiddleware())
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
// @Security ApiKeyAuth
// @Param device body DeviceCreate true "Device creation data"
// @Success 200 {object} Command
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var deviceCreate DeviceCreate
	if err := c.ShouldBindJSON(&deviceCreate); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	// Check if user has access to the home
	homeIDStr := deviceCreate.HomeID.String()
	userHomes, exists := c.Get("user_homes")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{Message: "User homes not found in context"})
		return
	}

	homes, ok := userHomes.([]string)
	if !ok {
		c.JSON(http.StatusUnauthorized, Response{Message: "Invalid user homes format"})
		return
	}

	// Check if user has access to the requested home
	hasAccess := false
	for _, home := range homes {
		if home == homeIDStr {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, Response{Message: "Access denied to this home"})
		return
	}

	command := Command{
		CommandType: CreateDevice,
		Params: map[string]string{
			"device_id":   uuid.New().String(),
			"device_type": string(deviceCreate.DeviceType),
			"name":        deviceCreate.Name,
			"home_id":     homeIDStr,
			"created_at":  time.Now().Format("2006-01-02T15:04:05Z07:00"),
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
// @Security ApiKeyAuth
// @Param device_id path string true "Device ID"
// @Param home_id query string true "Home ID"
// @Param command body DeviceAction true "Device command"
// @Success 200 {object} Command
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /devices/{device_id}/action [post]
func (h *DeviceHandler) SendAction(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	_, err := uuid.Parse(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "invalid device ID"})
		return
	}

	// Get home_id from query parameters
	homeIDStr := c.Query("home_id")
	if homeIDStr == "" {
		c.JSON(http.StatusBadRequest, Response{Message: "home_id query parameter is required"})
		return
	}

	// Validate home_id format
	_, err = uuid.Parse(homeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "invalid home_id format"})
		return
	}

	// Check if user has access to the home
	userHomes, exists := c.Get("user_homes")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{Message: "User homes not found in context"})
		return
	}

	homes, ok := userHomes.([]string)
	if !ok {
		c.JSON(http.StatusUnauthorized, Response{Message: "Invalid user homes format"})
		return
	}

	hasAccess := false
	for _, home := range homes {
		if home == homeIDStr {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, Response{Message: "Access denied to this home"})
		return
	}

	var action DeviceAction
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	params := map[string]string{
		"device_id":  deviceIDStr,
		"home_id":    homeIDStr,
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
