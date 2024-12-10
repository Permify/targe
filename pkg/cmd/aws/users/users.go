package users

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/config"
	"github.com/Permify/kivo/pkg/cmd/common"
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
func NewUsersCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "users",
		Short: "",
		RunE:  users(),
	}

	conf := config.DefaultConfig()
	f := command.Flags()
	f.StringP("config", "c", "", "config file (default is $HOME/.kivo.yaml)")

	f.String("user", conf.User, "user")
	f.String("action", conf.Action, "action")
	f.String("policy", conf.Policy, "policy")
	f.String("resource", conf.Resource, "resource")
	f.String("service", conf.Service, "service")
	f.String("policy-option", conf.PolicyOption, "policy option")

	// SilenceUsage is set to true to suppress usage when an error occurs
	command.SilenceUsage = true

	command.PreRun = func(cmd *cobra.Command, args []string) {
		RegisterUsersFlags(f)

		// Replace "requirements" with the actual path to your folder
		requirementsPath := "requirements"

		// Check if the requirements folder exists
		if folderExists(requirementsPath) {
			return
		}

		if _, err := tea.NewProgram(common.NewRequirements()).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}

	return command
}

func users() func(cmd *cobra.Command, args []string) error {
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

		var state State

		if cfg.User != "" {
			state.SetUser(&User{
				Name: cfg.User,
				Arn:  "arn:aws:iam::123456789012:user/" + cfg.User,
			})
		}

		if cfg.Action != "" {
			state.SetAction(&Action{
				Id:   cfg.Action,
				Name: ReachableActions[cfg.Action].Name,
				Desc: ReachableActions[cfg.Action].Desc,
			})
		}

		if cfg.Policy != "" {
			state.SetPolicy(&Policy{
				Name: cfg.Policy,
				Arn:  "arn:aws:iam::aws:policy/" + cfg.Policy,
			})
		}

		if cfg.Resource != "" {
			state.SetResource(&Resource{
				Name: cfg.Resource,
				Arn:  "arn:aws:iam::aws:resource/" + cfg.Resource,
			})
		}

		if cfg.Service != "" {
			state.SetService(&Service{
				Name: cfg.Service,
			})
		}

		if cfg.PolicyOption != "" {
			state.SetPolicyOption(&CustomPolicyOption{
				Name: cfg.PolicyOption,
				Desc: ReachableCustomPolicyOptions[cfg.PolicyOption].Desc,
			})
		}

		p := tea.NewProgram(RootModel(state.FindFlow()), tea.WithAltScreen())
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

func folderExists(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
