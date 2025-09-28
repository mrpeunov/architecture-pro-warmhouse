package providers

import (
	"bytes"
	"devices_command_consumer/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type WarmHouseProvider struct {
	BaseProvider
	baseURL    string
	httpClient *http.Client
}

type Sensor struct {
	ID       int     `json:"id,omitempty"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Location string  `json:"location"`
	Unit     string  `json:"unit"`
	Value    float64 `json:"value,omitempty"`
	Status   string  `json:"status,omitempty"`
}

type SensorCreate struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
	Unit     string `json:"unit"`
}

type SensorUpdate struct {
	Value  float64 `json:"value"`
	Status string  `json:"status"`
}

type SensorValueAction struct {
	DeviceID uuid.UUID `json:"device_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Action   string    `json:"action" binding:"required" example:"set_value"`
	Value    float64   `json:"value" binding:"required" example:"30"`
}

type SensorTurnOnAction struct {
	DeviceID uuid.UUID `json:"device_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Action   string    `json:"action" binding:"required" example:"turn_on"`
}

func NewWarmHouseProvider() *WarmHouseProvider {
	return &WarmHouseProvider{
		baseURL: getEnv("WARMHOUSE_API_URL", "http://localhost:8080"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (w *WarmHouseProvider) Create(device *models.Device) (string, error) {
	sensorCreate := &SensorCreate{
		Name:     device.Name,
		Type:     "temperature",
		Location: "Unknown",
		Unit:     "°C",
	}

	jsonData, err := json.Marshal(sensorCreate)
	if err != nil {
		return "", fmt.Errorf("failed to marshal sensor data: %w", err)
	}

	resp, err := w.httpClient.Post(
		w.baseURL+"/api/v1/sensors",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create sensor: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	var sensor Sensor
	err = json.Unmarshal(respBody, &sensor)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal sensor: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create sensor, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return string(sensor.ID), nil
}

func (w *WarmHouseProvider) SendAction(device *models.Device, data map[string]string) error {
	action := data["action"]

	sensorID, err := strconv.Atoi(device.OuterID)
	if err != nil {
		return fmt.Errorf("failed to convert sensor id to int: %w", err)
	}

	switch action {
	case "turn_on":
		update := &SensorUpdate{
			Value:  20.0,
			Status: "active",
		}
		return w.UpdateSensorValue(sensorID, update)

	case "turn_off":
		update := &SensorUpdate{
			Value:  0.0,
			Status: "inactive",
		}
		return w.UpdateSensorValue(sensorID, update)

	case "set_temperature":
		if tempStr, exists := data["value"]; exists {
			// Parse temperature value (simplified)
			var temp float64
			if _, err := fmt.Sscanf(tempStr, "%f", &temp); err == nil {
				update := &SensorUpdate{
					Value:  temp,
					Status: "active",
				}
				return w.UpdateSensorValue(sensorID, update)
			}
		}
		return fmt.Errorf("invalid temperature parameter")

	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}

func (w *WarmHouseProvider) UpdateSensorValue(sensorID int, update *SensorUpdate) error {
	jsonData, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal update data: %w", err)
	}

	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s/api/v1/sensors/%d/value", w.baseURL, sensorID),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update sensor value: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update sensor value, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
