// Package config loads and validates the service configuration from the
// environment, falling back to an optional local .env file.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultPort        = "8080"
	defaultDBHost      = "localhost"
	defaultDBPort      = "5432"
	defaultDBUser      = "postgres"
	defaultDBPassword  = "postgres"
	defaultDBName      = "eventer"
	defaultCORSOrigins = "http://localhost:5173"
)

// Config is the fully validated application configuration.
type Config struct {
	Port        string
	JWTSecret   string
	CORSOrigins []string
	Database    Database
}

// Database holds the PostgreSQL connection settings.
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// Load reads configuration from the environment (and an optional .env file)
// and fails if anything required is missing or malformed.
func Load() (*Config, error) {
	// A .env file is convenient locally but optional; ignore load errors.
	_ = godotenv.Load()

	cfg := &Config{
		Port:        env("PORT", defaultPort),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		CORSOrigins: splitList(env("CORS_ORIGINS", defaultCORSOrigins)),
		Database: Database{
			Host:     env("DB_HOST", defaultDBHost),
			Port:     env("DB_PORT", defaultDBPort),
			User:     env("DB_USER", defaultDBUser),
			Password: env("DB_PASSWORD", defaultDBPassword),
			Name:     env("DB_NAME", defaultDBName),
		},
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, fmt.Errorf("PORT must be a number, got %q", cfg.Port)
	}
	if cfg.Database.User == "" || cfg.Database.Name == "" {
		return nil, fmt.Errorf("DB_USER and DB_NAME are required")
	}
	if len(cfg.CORSOrigins) == 0 {
		return nil, fmt.Errorf("CORS_ORIGINS must list at least one origin")
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string for this database.
func (d Database) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		d.Host, d.User, d.Password, d.Name, d.Port,
	)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
