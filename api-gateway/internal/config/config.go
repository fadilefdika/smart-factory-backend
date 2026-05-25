package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	PostgresDSN     string
	InfluxDbUrl     string
	InfluxDbToken   string
	InfluxDBOrg     string
	InfluxDBBucket  string
	MqttBroker      string
	MQTTClientID    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, loading variables from system environment")
	}

	return &Config{
		AppPort:        getEnv("APP_PORT", "8000"),
		PostgresDSN:    getEnv("POSTGRES_DSN", "host=localhost user=admin password=adminpassword dbname=db_smartfactory_production port=5432 sslmode=disable"),
		InfluxDbUrl:    getEnv("INFLUXDB_URL", "http://localhost:8086"),
		InfluxDbToken:  getEnv("INFLUXDB_TOKEN", "super-secret-token-12345"),
		InfluxDBOrg:    getEnv("INFLUXDB_ORG", "smartfactory"),
		InfluxDBBucket: getEnv("INFLUXDB_BUCKET", "telemetry"),
		MqttBroker:     getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		MQTTClientID:   getEnv("MQTT_CLIENT_ID", "gateway-api-01"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
