package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	AppUrl     string
	AppPort    string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	DBSSLMode string

	LogLevel string
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, использую переменные окружения")
	}

	config := &Config{}

	config.DBHost = getEnv("DB_HOST", "localhost")

	dbPortStr := getEnv("DB_PORT", "5432")
	config.DBPort, err = strconv.Atoi(dbPortStr)
	if err != nil {
		panic(err)
	}

	config.AppEnv = getEnv("APP_ENV", "local")
	config.AppUrl = getEnv("APP_URL", "localhost")
	config.AppPort = getEnv("APP_PORT", "8080")
	config.DBUser = getEnv("DB_USER", "")
	config.DBPassword = getEnv("DB_PASSWORD", "")
	config.DBName = getEnv("DB_NAME", "")
	config.DBSSLMode = getEnv("DB_SSLMODE", "disable")

	config.LogLevel = getEnv("LOG_LEVEL", "info")

	if config.DBHost == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		panic("необходимые параметры базы данных отсутствуют")
	}

	return config
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
