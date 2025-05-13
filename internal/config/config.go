package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Environment represents the deployment environment
type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func (env Environment) IsDevelopment() bool {
	return env == Development
}

func (env Environment) IsProduction() bool {
	return env == Production
}

// Config represents the service configuration
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Logging LoggingConfig `mapstructure:"logging"`
}

// ServerConfig holds the server-specific configuration
type ServerConfig struct {
	HTTP HTTPConfig  `mapstructure:"http"`
	GRPC GRPCConfig  `mapstructure:"grpc"`
	Env  Environment `mapstructure:"environment"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=json text"`
}

// LoadConfig loads and validates the configuration
func LoadConfig(logger *logrus.Logger, relConfigPath string, env string) (*Config, error) {
	// Get the absolute path to the config file
	configPath, err := filepath.Abs(relConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for config file: %w", err)
	}

	// Get the directory of the config file
	configDir := filepath.Dir(configPath)
	configName := filepath.Base(configPath)
	configExt := filepath.Ext(configPath)
	configNameWithoutExt := configName[:len(configName)-len(configExt)]

	// Initialize viper
	v := viper.New()
	v.SetConfigName(configNameWithoutExt)
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	// Set defaults
	setDefaults(v)

	// Try to read environment-specific config
	envConfigPath := filepath.Join(configDir, fmt.Sprintf("%s.%s%s", configNameWithoutExt, env, configExt))
	if _, err := os.Stat(envConfigPath); err == nil {
		v.SetConfigName(fmt.Sprintf("%s.%s", configNameWithoutExt, env))
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to merge environment config: %w", err)
		}
	}

	logger.Infof("Using config file: %s", envConfigPath)

	// Read the default config
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// API defaults
	v.SetDefault("api.http_port", 8085)
	v.SetDefault("api.grpc_port", 8086)
	v.SetDefault("api.read_timeout", "5s")
	v.SetDefault("api.write_timeout", "5s")
	v.SetDefault("api.idle_timeout", "120s")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
}

// validateConfig validates the configuration using struct tags
func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				return fmt.Errorf("validation failed for field %s: %s", e.Field(), e.Tag())
			}
		}
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
