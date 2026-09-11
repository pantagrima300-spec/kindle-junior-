package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                  string
	ServerAddr              string
	PocketBaseURL           string
	PocketBaseAdminEmail    string
	PocketBaseAdminPassword string
	CORSOrigins             string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:                  getEnv("APP_ENV", "development"),
		ServerAddr:              getEnv("SERVER_ADDR", ":8080"),
		PocketBaseURL:           getEnv("POCKETBASE_URL", "http://127.0.0.1:8090"),
		PocketBaseAdminEmail:    getEnv("POCKETBASE_ADMIN_EMAIL", ""),
		PocketBaseAdminPassword: getEnv("POCKETBASE_ADMIN_PASSWORD", ""),
		CORSOrigins:             getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
