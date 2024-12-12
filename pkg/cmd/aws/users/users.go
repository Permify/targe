package users

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/config"
)

type Users struct {
	model tea.Model
}

func (m Users) Init() tea.Cmd {
	return m.model.Init() // rest methods are just wrappers for the model's methods
}

func (m Users) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m Users) View() string {
	return m.model.View()
}

// NewUsersCommand -
func NewUsersCommand(cfg *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "users",
		Short: "",
		RunE:  users(cfg),
	}

	f := command.Flags()

	f.String("user", "", "user")
	f.String("operation", "", "operation")
	f.String("policy", "", "policy")
	f.String("resource", "", "resource")
	f.String("service", "", "service")
	f.String("policy-option", "", "policy option")

	// SilenceUsage is set to true to suppress usage when an error occurs
	command.SilenceUsage = true

	command.PreRun = func(cmd *cobra.Command, args []string) {
		RegisterUsersFlags(f)
	}

	return command
}

func users(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var cfg *config.Config
		var err error
		cfgFile := viper.GetString("config.file")
		if cfgFile != "" {
			cfg, err = config.NewConfigWithFile(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to create new config: %w", err)
			}

			if err = viper.Unmarshal(cfg); err != nil {
				return fmt.Errorf("failed to unmarshal config: %w", err)
			}
		} else {
			// Load configuration
			cfg, err = config.NewConfig()
			if err != nil {
				return fmt.Errorf("failed to create new config: %w", err)
			}

			if err = viper.Unmarshal(cfg); err != nil {
				return fmt.Errorf("failed to unmarshal config: %w", err)
			}
		}

		// get min coverage from viper
		user := viper.GetString("user")
		operation := viper.GetString("operation")
		policy := viper.GetString("policy")
		resource := viper.GetString("resource")
		service := viper.GetString("service")
		policyOption := viper.GetString("policy-option")

		state, err := NewState(context.Background())
		if err != nil {
			return fmt.Errorf("failed to create new state: %w", err)
		}

		if user != "" {
			state.SetUser(&User{
				Name: user,
				Arn:  "arn:aws:iam::123456789012:user/" + user,
			})
		}

		if operation != "" {
			state.SetOperation(&Operation{
				Id:   operation,
				Name: ReachableOperations[operation].Name,
				Desc: ReachableOperations[operation].Desc,
			})
		}

		if policy != "" {
			state.SetPolicy(&Policy{
				Name: policy,
				Arn:  "arn:aws:iam::aws:policy/" + policy,
			})
		}

		if resource != "" {
			state.SetResource(&Resource{
				Name: resource,
				Arn:  "arn:aws:iam::aws:resource/" + resource,
			})
		}

		if service != "" {
			state.SetService(&Service{
				Name: service,
			})
		}

		if policyOption != "" {
			state.SetPolicyOption(&CustomPolicyOption{
				Name: policyOption,
				Desc: ReachableCustomPolicyOptions[policyOption].Desc,
			})
		}

		p := tea.NewProgram(RootModel(state.Next()), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		return nil
	}
}

func RootModel(m tea.Model) Users {
	return Users{
		model: m,
	}
}
