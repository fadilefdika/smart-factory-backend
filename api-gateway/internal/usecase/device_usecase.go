package usecase

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/smartfactory/api-gateway/internal/domain"
)

type deviceUsecase struct {
	deviceRepo domain.DeviceRepository
}

func NewDeviceUsecase(deviceRepo domain.DeviceRepository) domain.DeviceUsecase {
	return &deviceUsecase{
		deviceRepo: deviceRepo,
	}
}

func (u *deviceUsecase) RegisterDevice(device *domain.Device) error {
	device.ID = uuid.New()
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	device.LastSeen = time.Now()
	if device.Status == "" {
		device.Status = "inactive"
	}
	
	// TODO: Generate MQTT Credentials
	// In a real scenario, we would trigger an external service or generate a hashed password
	// to be added to mosquitto auth plugin.

	return u.deviceRepo.Create(device)
}

func (u *deviceUsecase) GetDeviceStatus(id uuid.UUID) (*domain.Device, error) {
	return u.deviceRepo.GetByID(id)
}

func (u *deviceUsecase) ListDevices(deviceType string) ([]domain.Device, error) {
	// For simplicity, we just list all devices.
	// If deviceType is needed, we would add a filter method to the repository.
	return u.deviceRepo.List()
}

func (u *deviceUsecase) SaveTelemetry(deviceID uuid.UUID, data string) error {
	var payload struct {
		Suhu float64 `json:"suhu"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err == nil {
		if payload.Suhu > 85.0 {
			resp, err := u.deviceRepo.TriggerAIInspection(deviceID)
			if err != nil {
				// Non-blocking failure: Hanya log error dan lanjutkan proses
				fmt.Printf("[Golang Usecase] ❌ Gagal memicu AI QC: %v\n", err)
			} else {
				// Berhasil parsing objek response
				fmt.Printf("[Golang Usecase] ⚠️ Overheat Terdeteksi (%.2f)! Memicu AI QC... Respon Python: %s (Confidence: %.2f%%)\n", 
					payload.Suhu, 
					resp.Result.Decision, 
					resp.Result.ConfidenceScore*100)
			}
		}
	}

	return u.deviceRepo.SaveTelemetry(deviceID, data)
}


func (u *deviceUsecase) GetTelemetryData()(map[uuid.UUID][]string, error){
	return u.deviceRepo.GetTelemetryData()
}