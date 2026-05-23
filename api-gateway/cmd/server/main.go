package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	
	"github.com/smartfactory/api-gateway/internal/config"
	deliveryHttp "github.com/smartfactory/api-gateway/internal/delivery/http"
	"github.com/smartfactory/api-gateway/internal/repository"
	"github.com/smartfactory/api-gateway/internal/usecase"
)

func main() {
	log.Println("Starting Smart Factory API Gateway...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize PostgreSQL Database Connection
	// db, err := sqlx.Connect("postgres", cfg.PostgresDSN)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	// }
	// defer db.Close()
	// log.Println("Connected to PostgreSQL successfully")

	// Initialize Repository & Usecase (In-Memory Local Mode)
	deviceRepo := repository.NewDeviceRepository()
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)

	// Initialize Gin HTTP Server
	r := gin.Default()
	
	// Register Handlers
	deliveryHttp.NewDeviceHandler(r, deviceUsecase)

	// TODO: Initialize MQTT Consumer
	// TODO: Initialize InfluxDB Client

	log.Printf("Starting HTTP server on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
