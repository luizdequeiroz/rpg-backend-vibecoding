package config

import (
	"os"
	"strconv"
	"time"
)

// Config contém todas as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Log      LogConfig
}

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	URL string
}

// AuthConfig configurações de autenticação
type AuthConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

// LogConfig configurações de log
type LogConfig struct {
	Level  string
	Format string
}

// Load carrega as configurações a partir das variáveis de ambiente
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("HOST", "localhost"),
			Port: getEnvAsInt("PORT", 8080),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
			JWTExpiration: getEnvAsDuration("JWT_EXPIRATION", "24h"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}
}

// getEnv obtém uma variável de ambiente com valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtém uma variável de ambiente como inteiro com valor padrão
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration obtém uma variável de ambiente como duration com valor padrão
func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
