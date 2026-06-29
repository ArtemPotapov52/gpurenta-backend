package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          int
	DatabaseURL   string
	JWTSecret     string
	GoogleID      string
	GoogleSecret  string
	FrontendURL   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using env vars", "error", err)
	} else {
		slog.Info("loaded .env file")
	}
	return &Config{
		Port:         getEnvInt("PORT", 8080),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gpurenta?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		GoogleID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
