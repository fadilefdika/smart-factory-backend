package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/smartfactory/api-gateway/internal/config"
	deliveryHttp "github.com/smartfactory/api-gateway/internal/delivery/http"
	"github.com/smartfactory/api-gateway/internal/domain"
	"github.com/smartfactory/api-gateway/internal/repository"
	"github.com/smartfactory/api-gateway/internal/usecase"
)

type TelemetryPayload struct {
	DeviceID  string  `json:"device_id"`
	Suhu      float64 `json:"suhu"`
	Getaran   float64 `json:"getaran"`
	Timestamp string  `json:"timestamp"`
}

func main() {
	log.Println("Starting Smart Factory API Gateway & Simulator...")

	// Load configuration
	cfg := config.LoadConfig()
	port := cfg.AppPort
	if port == "" {
		port = "8000"
	}

	// 1. Initialize In-Memory Repository & Usecase (ONLY ONCE)
	deviceRepo := repository.NewDeviceRepository()
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)

	// Buat Dummy Device agar tidak error "device not found" saat update LastSeen
	// Gunakan UUID yang statis/sama agar fetch dari Next.js konsisten
	dummyID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	deviceRepo.Create(&domain.Device{
		ID:        dummyID,
		Name:      "Mesin CNC 01 Simulator",
		Type:      "sensor",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	log.Printf("Created Dummy Device for simulation: %s\n", dummyID.String())

	// 2. Initialize Gin HTTP Server
	r := gin.Default()

	// Pasang middleware CORS resmi Gin
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// Register Handlers, bagikan instance usecase yang sama
	deliveryHttp.NewDeviceHandler(r, deviceUsecase)

	// 3. Buat Go Channel Internal untuk Simulator
	telemetryChan := make(chan string, 100)

	// 4. Jalankan Simulator (Background Goroutine)
	go internalSimulator(telemetryChan, dummyID)

	// 5. Logika Utama Consumer (Asynchronous Worker Pool)
	for i := 1; i <= 3; i++ {
		go func(workerID int) {
			for data := range telemetryChan {
				processTelemetryMessage(deviceUsecase, data, workerID)
			}
		}(i)
	}

	// Konfigurasi server HTTP
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 6. Jalankan Server HTTP (Background Goroutine agar tidak memblokir)
	go func() {
		log.Printf("Starting HTTP server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// 7. Proses Shutdown Graceful
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	// Tutup HTTP Server dengan timeout 5 detik
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("HTTP Server Shutdown:", err)
	}

	// Tutup channel telemetry Simulator
	close(telemetryChan)
	log.Println("Server and Consumer Simulator gracefully stopped")
}

// internalSimulator menghasilkan string JSON telemetri palsu setiap 1 detik
func internalSimulator(ch chan<- string, deviceID uuid.UUID) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		payload := TelemetryPayload{
			DeviceID:  deviceID.String(),
			Suhu:      60.0 + rng.Float64()*40.0, // 60.0 - 100.0
			Getaran:   0.5 + rng.Float64()*5.0,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		jsonData, err := json.Marshal(payload)
		if err == nil {
			// Cegah panic write on closed channel jika server dimatikan
			defer func() {
				recover()
			}()
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

	// Teruskan data ke Usecase (menggunakan instance memory yang sama)
	err = us.SaveTelemetry(devID, data)
	if err != nil {
		log.Printf("[Worker-%d] Gagal menyimpan telemetri: %v\n", workerID, err)
	} else {
		fmt.Printf("[Worker-%d] ✅ BERHASIL DIKONSUMSI & DISIMPAN: DeviceID=%s | Data=%s\n", workerID, payload.DeviceID, data)
	}
}
