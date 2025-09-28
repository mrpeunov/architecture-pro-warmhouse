package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"telemetry_api/models"
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

func NewWarmHouseProvider() *WarmHouseProvider {
	return &WarmHouseProvider{
		baseURL: getEnv("WARMHOUSE_API_URL", "http://localhost:8080"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (w *WarmHouseProvider) GetTelemetry(device *models.Device) (*models.TelemetryData, error) {
	sensorID, err := strconv.Atoi(device.OuterID)
	if err != nil {
		return nil, fmt.Errorf("invalid outer_id '%s': must be a valid integer", device.OuterID)
	}

	resp, err := w.httpClient.Get(fmt.Sprintf("%s/api/v1/sensors/%d", w.baseURL, sensorID))
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get sensor data, status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var sensor Sensor
	err = json.Unmarshal(body, &sensor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sensor: %w", err)
	}

	telemetry := &models.TelemetryData{
		TelemetryID: uuid.New(),
		DeviceID:    device.DeviceID,
		Value:       sensor.Value,
		CreatedAt:   time.Now(),
	}

	return telemetry, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
