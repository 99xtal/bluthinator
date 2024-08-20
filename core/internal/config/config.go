package config

import (
	"os"
)

type Config struct {
	PostgresHost          string
	PostgresPort          string
	PostgresUser          string
	PostgresPass          string
	PostgresDB            string
	ObjectStorageEndpoint string
}

func (c *Config) GetPostgresConnString() string {
	return "host=" + c.PostgresHost + " port=" + c.PostgresPort + " user=" + c.PostgresUser + " password=" + c.PostgresPass + " dbname=" + c.PostgresDB + " sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func New() *Config {
	return &Config{
		PostgresHost:          getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:          getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:          getEnv("POSTGRES_USER", "default_user"),
		PostgresPass:          getEnv("POSTGRES_PASSWORD", "default_pass"),
		PostgresDB:            getEnv("POSTGRES_DB", "default_db"),
		ObjectStorageEndpoint: getEnv("OBJECT_STORAGE_ENDPOINT", "http://localhost:9000"),
	}
}
