package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	LogLevel string

	ExternalAPIURL string

	ServerPort string
}

func MustLoad() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, использую переменные окружения")
	}

	config := &Config{}

	config.DBHost = getEnv("DB_HOST", "localhost")

	dbPortStr := getEnv("DB_PORT", "5432")
	config.DBPort, err = strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат DB_PORT: %v", err)
	}

	config.AppEnv = getEnv("APP_ENV", "dev")
	config.DBUser = getEnv("DB_USER", "")
	config.DBPassword = getEnv("DB_PASSWORD", "")
	config.DBName = getEnv("DB_NAME", "")
	config.DBSSLMode = getEnv("DB_SSLMODE", "")

	config.LogLevel = getEnv("LOG_LEVEL", "info")

	config.ExternalAPIURL = getEnv("EXTERNAL_API_URL", "")

	config.ServerPort = getEnv("PORT", "8080")

	if config.DBHost == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		return nil, fmt.Errorf("необходимые параметры базы данных отсутствуют")
	}

	return config, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
