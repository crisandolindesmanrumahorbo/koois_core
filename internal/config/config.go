package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	MaxConn  int
	MinConn  int
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		println("⚠️ No .env file found")
	}
	return &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", ""),
			MaxConn:  10,
			MinConn:  2,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
	}
}

func (c *Config) Validate() {
	missing := []string{}

	if c.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.Database.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.Database.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Database.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if c.JWT.Secret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "Missing required environment variables: %v\n", missing)
		os.Exit(1)
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
