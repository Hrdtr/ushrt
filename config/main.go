package config

import (
	"os"
)

func GetEnv(key string) string {
	value, _ := os.LookupEnv(key)
	return value
}

func GetEnvWithFallback(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok && fallback != "" {
		return fallback
	}
	return value
}

var (
	APP_ENV                = GetEnvWithFallback("APP_ENV", "development")
	APP_API_KEY            = GetEnvWithFallback("APP_API_KEY", "changeme")
	APP_BASE_URL           = GetEnvWithFallback("APP_BASE_URL", "http://localhost:3000")
	APP_CORS_ALLOW_ORIGINS = GetEnvWithFallback("APP_CORS_ALLOW_ORIGINS", "http://localhost:3000")

	POSTGRES_HOST     = GetEnvWithFallback("POSTGRES_HOST", "localhost")
	POSTGRES_PORT     = GetEnvWithFallback("POSTGRES_PORT", "5432")
	POSTGRES_USER     = GetEnvWithFallback("POSTGRES_USER", "pguser")
	POSTGRES_PASSWORD = GetEnvWithFallback("POSTGRES_PASSWORD", "pgpassword")
	POSTGRES_DB       = GetEnvWithFallback("POSTGRES_DB", "ushrt")
	POSTGRES_SSL_MODE = GetEnvWithFallback("POSTGRES_SSL_MODE", "disable")
)
