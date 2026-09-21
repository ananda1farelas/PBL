package config

import (
	"os"
	"time"
)

const (
	DefaultJWTSecret = "api-students-dev-secret-jangan-pakai-di-produksi"
	DefaultJWTIssuer = "api-students"
)

type Config struct {
	AppPort    string
	JWTSecret  string
	JWTIssuer  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func LoadConfig() *Config {
	LoadEnv()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	accessTTL := time.Duration(GetEnvInt("ACCESS_TOKEN_TTL_MINUTES", 15)) * time.Minute
	refreshTTL := time.Duration(GetEnvInt("REFRESH_TOKEN_TTL_DAYS", 30)) * 24 * time.Hour

	return &Config{
		AppPort:    port,
		JWTSecret:  GetEnv("JWT_SECRET", DefaultJWTSecret),
		JWTIssuer:  GetEnv("JWT_ISSUER", DefaultJWTIssuer),
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}
}