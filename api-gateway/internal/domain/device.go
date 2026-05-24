package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Type      string          `json:"type" db:"type"`       // sensor, actuator, gateway
	Status    string          `json:"status" db:"status"`   // active, inactive, offline
	Metadata  json.RawMessage `json:"metadata" db:"metadata"`
	LastSeen  time.Time       `json:"last_seen" db:"last_seen"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type DeviceRepository interface {
	Create(device *Device) error
	GetByID(id uuid.UUID) (*Device, error)
	UpdateStatus(id uuid.UUID, status string) error
	List() ([]Device, error)
	SaveTelemetry(deviceID uuid.UUID, data string) error
	TriggerAIInspection(deviceID uuid.UUID) (*AIInspectionResponse, error)
	GetTelemetryData()(map[uuid.UUID][]string, error)
}

type AIInspectionResult struct {
	ID              string  `json:"id"`
	ImageURL        string  `json:"image_url"`
	Decision        string  `json:"decision"`
	ConfidenceScore float64 `json:"confidence_score"`
	InspectedAt     string  `json:"inspected_at"`
}

type AIInspectionResponse struct {
	Message string             `json:"message"`
	Result  AIInspectionResult `json:"result"`
}

type DeviceUsecase interface {
	RegisterDevice(device *Device) error
	GetDeviceStatus(id uuid.UUID) (*Device, error)
	ListDevices(deviceType string) ([]Device, error)
	SaveTelemetry(deviceID uuid.UUID, data string) error
	GetTelemetryData()(map[uuid.UUID][]string,error)
}
