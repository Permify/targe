package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Permify/targe/internal/config"
)

// NewConfigCommand - returns a new cobra command for config
func NewConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Manage targe configuration",
	}

	// Add subcommands
	command.AddCommand(newConfigSetCommand())
	command.AddCommand(newConfigGetCommand())

	return command
}

// newConfigSetCommand - returns a cobra command for setting config
func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration key-value pair",
		Args:  cobra.ExactArgs(2), // Requires exactly 2 arguments: key and value
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			if key != "openai_api_key" {
				return fmt.Errorf("invalid key: %s", key)
			}

			// Start with the default configuration values
			cfg := config.DefaultConfig()

			// Set the name and type of the config file to be read
			viper.SetConfigName("config")
			viper.SetConfigType("toml")

			// Add the path where the config file is located
			configPath := os.ExpandEnv("$HOME/.targe/")
			viper.AddConfigPath(configPath)

			// Ensure the directory exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				if err := os.MkdirAll(configPath, 0o755); err != nil {
					return fmt.Errorf("failed to create config directory: %w", err)
				}
			}

			// Update the key-value pair
			cfg.OpenaiApiKey = value

			// Read the config file
			err := viper.ReadInConfig()
			if err != nil {
				// If the error is due to the file not being found, create a new one
				var configFileNotFoundError viper.ConfigFileNotFoundError
				if errors.As(err, &configFileNotFoundError) {
					filePath := filepath.Join(configPath, "config.toml")
					if err := writeConfig(filePath, cfg); err != nil {
						return fmt.Errorf("failed to create config file: %w", err)
					}
				}
			}

			filePath := filepath.Join(configPath, "config.toml")
			if err := writeConfig(filePath, cfg); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}

			fmt.Printf("Configuration set: %s=%s\n", key, value)
			return nil
		},
	}
}

// newConfigGetCommand - returns a cobra command for getting config values
func newConfigGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value by key",
		Args:  cobra.ExactArgs(1), // Requires exactly 1 argument: key
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			if key != "openai_api_key" {
				return fmt.Errorf("invalid key: %s", key)
			}

			// Start with the default configuration values
			cfg := config.DefaultConfig()

			// Set the name and type of the config file to be read
			viper.SetConfigName("config")
			viper.SetConfigType("toml")

			// Add the path where the config file is located
			configPath := os.ExpandEnv("$HOME/.targe/")
			viper.AddConfigPath(configPath)

			// Ensure the directory exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				if err := os.MkdirAll(configPath, 0o755); err != nil {
					return fmt.Errorf("failed to create config directory: %w", err)
				}
			}

			// Read the config file
			err := viper.ReadInConfig()
			if err != nil {
				// If the error is due to the file not being found, create a new one
				var configFileNotFoundError viper.ConfigFileNotFoundError
				if errors.As(err, &configFileNotFoundError) {
					filePath := filepath.Join(configPath, "config.toml")
					if err := writeConfig(filePath, cfg); err != nil {
						return fmt.Errorf("failed to create config file: %w", err)
					}
				}
			}

			// Unmarshal the configuration data into the Config struct
			if err = viper.Unmarshal(cfg); err != nil {
				// If there's an error during unmarshalling, return the error with a message
				return fmt.Errorf("failed to unmarshal server config: %w", err)
			}

			fmt.Printf("%s=%s\n", key, cfg.OpenaiApiKey)
			return nil
		},
	}
}

func writeConfig(filePath string, cfg *config.Config) error {
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
