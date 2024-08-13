package config

// Config holds configuration for the application
type Config struct {
	DatabaseDSN string // Database Data Source Name
}

// LoadConfig loads configuration from a file or environment variables
func LoadConfig() (*Config, error) {
	// Implement your configuration loading logic here
	// For example, you can load from a file or environment variables
	return &Config{
		DatabaseDSN: "user:password@tcp(localhost:3306)/dbname",
	}, nil
}
