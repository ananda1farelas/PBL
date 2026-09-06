package config

import "os"

type Config struct {
	AppPort string
}

func LoadConfig() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	return &Config{
		AppPort: port,
	}
}
