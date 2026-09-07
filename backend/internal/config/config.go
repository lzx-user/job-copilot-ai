package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	AppEnv          string
	FrontendOrigin  string
	SupabaseURL     string
	SupabaseAnonKey string
	AIBaseURL       string
	AIAPIKey        string
	AIModel         string
	AITimeout       time.Duration
}

func Load() Config {
	// .env 只用于本地开发，因此文件不存在时继续读取系统环境变量。
	_ = godotenv.Load()

	return Config{
		AppPort:         getEnv("APP_PORT", "8080"),
		AppEnv:          getEnv("APP_ENV", "development"),
		FrontendOrigin:  getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey: getEnv("SUPABASE_ANON_KEY", ""),
		AIBaseURL:       getEnv("AI_API_BASE_URL", ""),
		AIAPIKey:        getEnv("AI_API_KEY", ""),
		AIModel:         getEnv("AI_MODEL", ""),
		AITimeout:       time.Duration(getEnvInt("AI_TIMEOUT_SECONDS", 30)) * time.Second,
	}
}

func getEnvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
