package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	FrontendOrigin string
}

func Load() Config {
	// .env 只用于本地开发，因此文件不存在时继续读取系统环境变量。
	_ = godotenv.Load()

	return Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
