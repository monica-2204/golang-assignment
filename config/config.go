package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds configuration values
type Config struct {
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	JWTSecret        string
	ServerPort       string
	LogLevel         string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default or environment variables")
	}

	cfg := &Config{
		DatabaseUser:     getEnv("DATABASE_USER", "root"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "Monu@2002"),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "3306"),
		DatabaseName:     getEnv("DATABASE_NAME", "student"),
		JWTSecret:        getEnv("JWT_SECRET", "3x@mP1e$eCr3t!VeRy$l0Ng@p@$sw0Rd"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

// getEnv retrieves the value of an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
