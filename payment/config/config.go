package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	WebhookSecret  string
	AcquiremockURL string
	BaseURL        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = buildPostgresURLFromEnv()
	}

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    dbURL,
		WebhookSecret:  mustEnv("WEBHOOK_SECRET"),
		AcquiremockURL: getEnv("ACQUIREMOCK_URL", "http://localhost:8000"),
		BaseURL:        mustEnv("BASE_URL"),
	}
	return cfg, nil
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
		panic(fmt.Sprintf("env %s is required", key))
	}
	return v
}

func buildPostgresURLFromEnv() string {
	db := getEnv("POSTGRES_DB", "postgres")
	user := getEnv("POSTGRES_USER", "postgres")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	sslmode := getEnv("POSTGRES_SSLMODE", "disable")

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   db,
	}
	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	return u.String()
}
