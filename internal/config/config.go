package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by Viper from a config file and/or environment variables.
type Config struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

// New loads configuration from file and environment variables.
func New() (*Config, error) {
	// --- Set up Viper ---

	// 1. Set default values (lowest priority)
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://boarded:lN9j11+P+99m)frFCF@localhost:5432/boarded_db?sslmode=disable")
	viper.SetDefault("JWT_SECRET", "a-very-secret-and-long-key-that-should-be-changed")
	viper.SetDefault("LOG_LEVEL", "debug")

	// 2. Read from a config file (e.g., config.yaml)
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // Look for config in the current directory

	// Attempt to read the config file. Ignore errors if the file doesn't exist.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// 3. Read from environment variables (highest priority)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// --- Unmarshal config into our struct ---
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}
