package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	TokenSecret string
	TokenTTL    time.Duration
	StorageURL  string
	PaymentURL  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	ttl, err := time.ParseDuration(getEnv("TOKEN_TTL", "720h"))
	if err != nil {
		return nil, fmt.Errorf("invalid TOKEN_TTL: %w", err)
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: mustEnv("DATABASE_URL"),
		TokenSecret: mustEnv("TOKEN_SECRET"),
		TokenTTL:    ttl,
		StorageURL:  mustEnv("STORAGE_URL"),
		PaymentURL:  mustEnv("PAYMENT_URL"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env %q is not set", key))
	}
	return v
}
