// Package config loads application settings from environment variables.
package config

import (
	"fmt"
	"net/url"
	"os"
)

// Config aggregates all subsystem configurations.
type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
	MinIO    MinIOConfig
}

// HTTPConfig holds the server listen address.
type HTTPConfig struct {
	Addr string
}

// PostgresConfig holds the DSN for pgx.
type PostgresConfig struct {
	DSN string
}

// MinIOConfig holds MinIO connection settings.
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

// Load reads config from env with sensible local defaults.
func Load() *Config {
	dbUrl := getenv("POSTGRES_DSN", buildPostgresURLFromEnv())

	return &Config{
		HTTP: HTTPConfig{
			Addr: getenv("HTTP_ADDR", ":8080"),
		},
		Postgres: PostgresConfig{
			DSN: dbUrl,
		},
		MinIO: MinIOConfig{
			Endpoint:  getenv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getenv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getenv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:    false,
			Bucket:    getenv("MINIO_BUCKET", "files"),
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func buildPostgresURLFromEnv() string {
	db := getenv("POSTGRES_DB", "postgres")
	user := getenv("POSTGRES_USER", "postgres")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := getenv("POSTGRES_HOST", "localhost")
	port := getenv("POSTGRES_PORT", "5432")
	sslmode := getenv("POSTGRES_SSLMODE", "disable")

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
