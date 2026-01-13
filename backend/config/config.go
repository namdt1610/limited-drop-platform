package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL string
	Port        string
	Environment string
	// Database pool settings
	MaxWriteConns int
	MaxReadConns  int
	BusyTimeout   int // milliseconds

	// AWS/LocalStack Configuration
	AWS AWSConfig
}

// AWSConfig holds AWS-specific settings
type AWSConfig struct {
	Endpoint   string
	Region     string
	AccessKey  string
	SecretKey  string
	UseSQS     bool
	UseS3      bool
	UseSecrets bool
}

// Load loads configuration from environment variables
func Load() *Config {
	config := &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "./database.db"),
		Port:          getEnv("PORT", "3030"),
		Environment:   getEnv("ENV", "development"),
		MaxWriteConns: getEnvAsInt("MAX_WRITE_CONNS", 1),    // SQLite writer uses 1 connection
		MaxReadConns:  getEnvAsInt("MAX_READ_CONNS", 100),   // Reader supports 100 concurrent connections
		BusyTimeout:   getEnvAsInt("DB_BUSY_TIMEOUT", 5000), // 5 seconds

		AWS: AWSConfig{
			Endpoint:   getEnv("AWS_ENDPOINT_URL", ""),
			Region:     getEnv("AWS_DEFAULT_REGION", "us-east-1"),
			AccessKey:  getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretKey:  getEnv("AWS_SECRET_ACCESS_KEY", ""),
			UseSQS:     getEnv("USE_SQS", "false") == "true",
			UseS3:      getEnv("USE_S3", "false") == "true",
			UseSecrets: getEnv("USE_SECRETS_MANAGER", "false") == "true",
		},
	}

	// Handle legacy DB_PATH if DATABASE_URL not set
	if config.DatabaseURL == "" {
		config.DatabaseURL = getEnv("DB_PATH", "./database.db")
	}

	return config
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
