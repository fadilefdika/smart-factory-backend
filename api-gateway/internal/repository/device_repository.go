package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/smartfactory/api-gateway/internal/domain"
)

type deviceRepository struct {
	mu        sync.RWMutex
	devices   map[uuid.UUID]*domain.Device
	telemetry map[uuid.UUID][]string
}

// In-Memory Repository
func NewDeviceRepository() domain.DeviceRepository {
	return &deviceRepository{
		devices:   make(map[uuid.UUID]*domain.Device),
		telemetry: make(map[uuid.UUID][]string),
	}
}

func (r *deviceRepository) Create(device *domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices[device.ID] = device
	return nil
}

func (r *deviceRepository) GetByID(id uuid.UUID) (*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	device, exists := r.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}
	return device, nil
}

func (r *deviceRepository) UpdateStatus(id uuid.UUID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	device, exists := r.devices[id]
	if !exists {
		return errors.New("device not found")
	}
	device.Status = status
	device.UpdatedAt = time.Now()
	return nil
}

func (r *deviceRepository) List() ([]domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Device
	for _, device := range r.devices {
		list = append(list, *device)
	}
	return list, nil
}

func (r *deviceRepository) SaveTelemetry(deviceID uuid.UUID, data string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Pastikan map untuk device ini ada
	if _, exists := r.telemetry[deviceID]; !exists {
		r.telemetry[deviceID] = make([]string, 0)
	}
	
	r.telemetry[deviceID] = append(r.telemetry[deviceID], data)
	
	// Update LastSeen perangkat
	if device, exists := r.devices[deviceID]; exists {
		device.LastSeen = time.Now()
	}
	
	return nil
}

func (r *deviceRepository) GetTelemetryData() (map[uuid.UUID][]string, error){
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[uuid.UUID][]string)

	for deviceID,message := range r.telemetry{
		result[deviceID] = message
	}

	return result, nil
} 

func (r *deviceRepository) TriggerAIInspection(deviceID uuid.UUID) (*domain.AIInspectionResponse, error) {
	// 1. Ambil URL AI Service dari environment, default ke localhost:8000
	aiServiceURL := os.Getenv("AI_QC_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8000"
	}

	// 2. Buat buffer untuk menampung body form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Buat dummy file (teks acak sebagai simulasi gambar)
	fw, err := w.CreateFormFile("image", "trigger_overheat.jpg")
	if err != nil {
		return nil, fmt.Errorf("gagal membuat form file image: %w", err)
	}
	_, err = io.WriteString(fw, "dummy image content for simulation")
	if err != nil {
		return nil, fmt.Errorf("gagal menulis ke form file: %w", err)
	}
	w.Close()

	// 3. Lakukan HTTP POST ke layanan Python FastAPI
	endpoint := fmt.Sprintf("%s/api/v1/qc/inspect", aiServiceURL)
	req, err := http.NewRequest("POST", endpoint, &b)
	if err != nil {
		return nil, fmt.Errorf("gagal menyusun request HTTP: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengeksekusi request ke AI service (koneksi mati?): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mendapat HTTP %d dari AI service", resp.StatusCode)
	}

	// 4. Parsing response ke struct AIInspectionResponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response body: %w", err)
	}

	var inspectionResp domain.AIInspectionResponse
	if err := json.Unmarshal(respBody, &inspectionResp); err != nil {
		return nil, fmt.Errorf("gagal mem-parsing JSON response: %w", err)
	}

	return &inspectionResp, nil
}
