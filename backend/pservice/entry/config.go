package main

import (
	"fmt"
	"os"
)

type Config struct {
	BackendHost string
	BackendPort string

	DbUser     string
	DbPassword string
	DbName     string
	DbHost     string
	DbPort     string

	JwtSecret string

	CorsAllowedOrigins string
}

func (c *Config) ConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName,
	)
}

func MustLoadConfig() *Config {
	cfg := &Config{}

	cfg.BackendHost = os.Getenv("BACKEND_HOST")
	cfg.BackendPort = optionalEnv("BACKEND_PORT", "8080")

	cfg.DbUser = requireEnv("DB_USER")
	cfg.DbPassword = requireEnv("DB_PASSWORD")
	cfg.DbName = requireEnv("DB_NAME")
	cfg.DbHost = requireEnv("DB_HOST")
	cfg.DbPort = requireEnv("DB_PORT")

	cfg.JwtSecret = requireEnv("JWT_SECRET")
	if len(cfg.JwtSecret) < 32 {
		panic("JWT_SECRET must be at least 32 characters")
	}

	val, exists := os.LookupEnv("CORS_ALLOWED_ORIGINS")
	if !exists {
		panic("CORS_ALLOWED_ORIGINS environment variable is required; use empty string for same-origin only")
	}
	cfg.CorsAllowedOrigins = val

	return cfg
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return val
}

func optionalEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
