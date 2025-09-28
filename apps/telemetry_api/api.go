package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"telemetry_api/models"
)

type ApiHandler struct {
	deviceService    *DeviceService
	telemetryService *TelemetryService
}

func NewApiHandler(deviceService *DeviceService, telemetryService *TelemetryService) *ApiHandler {
	return &ApiHandler{
		deviceService:    deviceService,
		telemetryService: telemetryService,
	}
}

func (h *ApiHandler) RegisterRoutes(router *gin.RouterGroup) {
	devices := router.Group("/devices")
	{
		devices.GET("", h.GetDevicesByHomeID)
		devices.GET("/:device_id/telemetry", h.GetLatestTelemetryByDeviceID)
	}
}

// GetDevicesByHomeID retrieves all devices for a specific home
// @Summary Get devices by home ID
// @Description Retrieve all devices for a specific home
// @Tags devices
// @Accept json
// @Produce json
// @Param home_id query string true "Home ID"
// @Success 200 {array} models.DeviceResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /devices [get]
func (h *ApiHandler) GetDevicesByHomeID(c *gin.Context) {
	homeIDStr := c.Query("home_id")
	homeID, err := uuid.Parse(homeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid home ID",
			Message: "Home ID must be a valid UUID",
		})
		return
	}

	devices, err := h.deviceService.GetDevicesByHomeID(homeID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Devices not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, devices)
}

// GetLatestTelemetryByDeviceID retrieves the latest telemetry data for a specific device
// @Summary Get latest telemetry data by device ID
// @Description Retrieve the most recent telemetry data for a specific device
// @Tags telemetry
// @Accept json
// @Produce json
// @Param device_id path string true "Device ID"
// @Success 200 {object} models.TelemetryResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /devices/{device_id}/telemetry [get]
func (h *ApiHandler) GetLatestTelemetryByDeviceID(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid device ID",
			Message: "Device ID must be a valid UUID",
		})
		return
	}

	telemetry, err := h.telemetryService.GetLatestTelemetryByDeviceID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Telemetry not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, telemetry)
}
