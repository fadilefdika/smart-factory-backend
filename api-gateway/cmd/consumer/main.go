package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	
	"github.com/smartfactory/api-gateway/internal/repository"
	"github.com/smartfactory/api-gateway/internal/usecase"
	"github.com/smartfactory/api-gateway/internal/domain"
)

type TelemetryPayload struct {
	DeviceID  string  `json:"device_id"`
	Suhu      float64 `json:"suhu"`
	Getaran   float64 `json:"getaran"`
	Timestamp string  `json:"timestamp"`
}

func main() {
	log.Println("Starting Smart Factory MQTT Consumer (In-Memory Simulator Mode)...")

	// 1. Initialize In-Memory Repository & Usecase
	deviceRepo := repository.NewDeviceRepository()
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)

	// Buat Dummy Device agar tidak error "device not found" saat update LastSeen
	dummyID := uuid.New()
	deviceRepo.Create(&domain.Device{
		ID:        dummyID,
		Name:      "Mesin CNC 01 Simulator",
		Type:      "sensor",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	log.Printf("Created Dummy Device for simulation: %s\n", dummyID.String())

	// 2. Buat Go Channel Internal
	telemetryChan := make(chan string, 100)

	// 3. Jalankan Simulator (Background Goroutine)
	go internalSimulator(telemetryChan, dummyID)

	// 4. Proses Shutdown Graceful
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 5. Logika Utama Consumer (Asynchronous Worker Pool)
	// Kita bisa menjalankan beberapa worker untuk mensimulasikan concurrent processing
	for i := 1; i <= 3; i++ {
		go func(workerID int) {
			for data := range telemetryChan {
				processTelemetryMessage(deviceUsecase, data, workerID)
			}
		}(i)
	}

	<-quit
	log.Println("Shutting down Consumer Simulator...")
	close(telemetryChan)
	log.Println("Consumer Simulator gracefully stopped")
}

// internalSimulator menghasilkan string JSON telemetri palsu setiap 1 detik
func internalSimulator(ch chan<- string, deviceID uuid.UUID) {
	rand.Seed(time.Now().UnixNano())
	for {
		payload := TelemetryPayload{
			DeviceID:  deviceID.String(),
			Suhu:      60.0 + rand.Float64()*40.0, // 60.0 - 100.0
			Getaran:   0.5 + rand.Float64()*5.0,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		
		jsonData, err := json.Marshal(payload)
		if err == nil {
			ch <- string(jsonData)
		}
		
		time.Sleep(1 * time.Second)
	}
}

// processTelemetryMessage meneruskan data ke Usecase hingga tersimpan di map lokal
func processTelemetryMessage(us domain.DeviceUsecase, data string, workerID int) {
	var payload TelemetryPayload
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Printf("[Worker-%d] Failed to parse JSON: %v", workerID, err)
		return
	}
	
	devID, err := uuid.Parse(payload.DeviceID)
	if err != nil {
		log.Printf("[Worker-%d] Invalid UUID format: %s", workerID, payload.DeviceID)
		return
	}

	// Teruskan data ke Usecase
	err = us.SaveTelemetry(devID, data)
	if err != nil {
		log.Printf("[Worker-%d] Gagal menyimpan telemetri: %v\n", workerID, err)
	} else {
		fmt.Printf("[Worker-%d] ✅ BERHASIL DIKONSUMSI & DISIMPAN: DeviceID=%s | Data=%s\n", workerID, payload.DeviceID, data)
	}
}
