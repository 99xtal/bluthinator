package config

import (
	"os"

	"github.com/elastic/go-elasticsearch/v7"
)

type Config struct {
	PostgresHost          string
	PostgresPort          string
	PostgresUser          string
	PostgresPass          string
	PostgresDB            string
	ElasticHost           string
	ElasticPort           string
	ElasticUser           string
	ElasticPass           string
	ObjectStorageEndpoint string
	ObjectStorageUser     string
	ObjectStoragePass     string
}

func (c *Config) GetPostgresConnString() string {
	return "host=" + c.PostgresHost + " port=" + c.PostgresPort + " user=" + c.PostgresUser + " password=" + c.PostgresPass + " dbname=" + c.PostgresDB + " sslmode=disable"
}

func (c *Config) GetElasticSearchConfig() elasticsearch.Config {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://" + c.ElasticHost + ":" + c.ElasticPort,
		},
		Username: c.ElasticUser,
		Password: c.ElasticPass,
	}
	return cfg
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
		ElasticHost:           getEnv("ELASTIC_HOST", "localhost"),
		ElasticPort:           getEnv("ELASTIC_PORT", "9200"),
		ElasticUser:           getEnv("ELASTIC_USER", "default_user"),
		ElasticPass:           getEnv("ELASTIC_PASS", "default_pass"),
		ObjectStorageEndpoint: getEnv("OBJECT_STORAGE_ENDPOINT", "http://localhost:9000"),
		ObjectStorageUser:     getEnv("OBJECT_STORAGE_USER", "minio"),
		ObjectStoragePass:     getEnv("OBJECT_STORAGE_PASSWORD", "minio123"),
	}
}
