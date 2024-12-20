package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type (
	Config struct {
		OpenaiApiKey string `mapstructure:"openai_api_key"`
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
	viper.SetConfigType("toml")

	// Add the path where the config file is located
	configPath := os.ExpandEnv("$HOME/.kivo/")
	viper.AddConfigPath(configPath)

	// Ensure the directory exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		// If the error is due to the file not being found, create a new one
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			filePath := filepath.Join(configPath, "config.toml")
			if err := writeDefaultConfig(filePath, cfg); err != nil {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
		}
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

	err := isTOML(dir)
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

func writeDefaultConfig(filePath string, cfg *Config) error {
	// Use viper to write the default configuration to a file
	viper.Set("openai_api_key", cfg.OpenaiApiKey)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := viper.WriteConfigAs(filePath); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}
	return nil
}

// DefaultConfig - Creates default config.
func DefaultConfig() *Config {
	return &Config{
		OpenaiApiKey: "",
	}
}

func isTOML(file string) error {
	ext := filepath.Ext(file)
	if ext != ".toml" {
		return errors.New("file is not toml")
	}
	return nil
}
