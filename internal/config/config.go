package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type (
	Config struct {
		User         string `mapstructure:"user"`
		Action       string `mapstructure:"action"`
		Policy       string `mapstructure:"policy"`
		PolicyOption string `mapstructure:"policy_option"`
		Service      string `mapstructure:"service"`
		Resource     string `mapstructure:"resource"`
	}
)

// NewConfig initializes and returns a new Config object by reading and unmarshalling
// the configuration file from the given path. It falls back to the DefaultConfig if the
// file is not found. If there's an error during the process, it returns the error.
func NewConfig() (*Config, error) {
	// Start with the default configuration values
	cfg := DefaultConfig()

	// Set the name and type of the config file to be read
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add the path where the config file is located
	viper.AddConfigPath("./config")

	// Read the config file
	err := viper.ReadInConfig()
	// If there's an error during reading the config file
	if err != nil {
		// Check if the error is because of the config file not being found
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); !ok {
			// If it's not a "file not found" error, return the error with a message
			return nil, fmt.Errorf("failed to load server config: %w", err)
		}
		// If it's a "file not found" error, the code will continue and use the default config
	}

	// Unmarshal the configuration data into the Config struct
	if err = viper.Unmarshal(cfg); err != nil {
		// If there's an error during unmarshalling, return the error with a message
		return nil, fmt.Errorf("failed to unmarshal server config: %w", err)
	}

	// Return the populated Config object
	return cfg, nil
}

// NewConfigWithFile initializes and returns a new Config object by reading and unmarshalling
// the configuration file from the given path. It falls back to the DefaultConfig if the
// file is not found. If there's an error during the process, it returns the error.
func NewConfigWithFile(dir string) (*Config, error) {
	// Start with the default configuration values
	cfg := DefaultConfig()

	viper.SetConfigFile(dir)

	err := isYAML(dir)
	if err != nil {
		return nil, err
	}

	// Read the config file
	err = viper.ReadInConfig()
	// If there's an error during reading the config file
	if err != nil {
		// Check if the error is because of the config file not being found
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); !ok {
			// If it's not a "file not found" error, return the error with a message
			return nil, fmt.Errorf("failed to load server config: %w", err)
		}
		if ok := errors.As(err, &viper.ConfigMarshalError{}); !ok {
			// If it's not a "file not found" error, return the error with a message
			return nil, fmt.Errorf("failed to load server config: %w", err)
		}
		// If it's a "file not found" error, the code will continue and use the default config
	}

	// Unmarshal the configuration data into the Config struct
	if err = viper.Unmarshal(cfg); err != nil {
		// If there's an error during unmarshalling, return the error with a message
		return nil, fmt.Errorf("failed to unmarshal server config: %w", err)
	}

	// Return the populated Config object
	return cfg, nil
}

// DefaultConfig - Creates default config.
func DefaultConfig() *Config {
	return &Config{
		User:         "",
		Action:       "",
		Policy:       "",
		PolicyOption: "",
		Service:      "",
		Resource:     "",
	}
}

func isYAML(file string) error {
	ext := filepath.Ext(file)
	if ext != ".yaml" {
		return errors.New("file is not yaml")
	}
	return nil
}
