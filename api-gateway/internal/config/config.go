package config

import (
	"os"
)

type Config struct {
	AppPort         string
	PostgresDSN     string
	InfluxDBURL     string
	InfluxDBToken   string
	InfluxDBOrg     string
	InfluxDBBucket  string
	MQTTBroker      string
	MQTTClientID    string
}

func LoadConfig() *Config {
	return &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		PostgresDSN:    getEnv("POSTGRES_DSN", "host=localhost user=admin password=adminpassword dbname=db_smartfactory_production port=5432 sslmode=disable"),
		InfluxDBURL:    getEnv("INFLUXDB_URL", "http://localhost:8086"),
		InfluxDBToken:  getEnv("INFLUXDB_TOKEN", "super-secret-token-12345"),
		InfluxDBOrg:    getEnv("INFLUXDB_ORG", "smartfactory"),
		InfluxDBBucket: getEnv("INFLUXDB_BUCKET", "telemetry"),
		MQTTBroker:     getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		MQTTClientID:   getEnv("MQTT_CLIENT_ID", "gateway-api-01"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
