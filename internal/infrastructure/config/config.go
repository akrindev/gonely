package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Log      LogConfig
	CORS     CORSConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name     string
	Env      string
	URL      string
	Port     string
	BasePath string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string
	JWTExpiry     time.Duration
	SessionExpiry time.Duration
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Try to read .env file, but don't fail if it doesn't exist
	_ = viper.ReadInConfig()

	// Set defaults
	setDefaults()

	config := &Config{
		App: AppConfig{
			Name:     viper.GetString("APP_NAME"),
			Env:      viper.GetString("APP_ENV"),
			URL:      viper.GetString("APP_URL"),
			Port:     viper.GetString("APP_PORT"),
			BasePath: viper.GetString("BASE_PATH"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
			Timezone: viper.GetString("DB_TIMEZONE"),
		},
		Auth: AuthConfig{
			JWTSecret:     viper.GetString("JWT_SECRET"),
			JWTExpiry:     viper.GetDuration("JWT_EXPIRY"),
			SessionExpiry: viper.GetDuration("SESSION_EXPIRY"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods: viper.GetStringSlice("CORS_ALLOWED_METHODS"),
			AllowedHeaders: viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	viper.SetDefault("APP_NAME", "Gonely")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_URL", "http://localhost:8080")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("BASE_PATH", "/api")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "salju")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "bunely")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_TIMEZONE", "Asia/Jakarta")

	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_EXPIRY", "168h")
	viper.SetDefault("SESSION_EXPIRY", "168h")

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"})
	viper.SetDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"})
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Auth.JWTSecret == "" || len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET is required and must be at least 32 characters")
	}
	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.Timezone,
	)
}

// IsDevelopment returns true if the app is in development mode
func (c *AppConfig) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if the app is in production mode
func (c *AppConfig) IsProduction() bool {
	return c.Env == "production"
}
