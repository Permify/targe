package groups

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Permify/kivo/internal/config"
)

type Groups struct {
	model tea.Model
}

func (m Groups) Init() tea.Cmd {
	return m.model.Init() // rest methods are just wrappers for the model's methods
}

func (m Groups) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m Groups) View() string {
	return m.model.View()
}

// NewGroupsCommand -
func NewGroupsCommand(cfg *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "groups",
		Short: "",
		RunE:  groups(cfg),
	}

	f := command.Flags()

	f.String("group", "", "group")
	f.String("operation", "", "operation")
	f.String("policy", "", "policy")
	f.String("resource", "", "resource")
	f.String("service", "", "service")
	f.String("policy-option", "", "policy option")

	// SilenceUsage is set to true to suppress usage when an error occurs
	command.SilenceUsage = true

	command.PreRun = func(cmd *cobra.Command, args []string) {
		RegisterGroupsFlags(f)
	}

	return command
}

func groups(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// get min coverage from viper
		group := viper.GetString("group")
		operation := viper.GetString("operation")
		policy := viper.GetString("policy")
		resource := viper.GetString("resource")
		service := viper.GetString("service")
		policyOption := viper.GetString("policy-option")

		state := &State{}
		if group != "" {
			state.SetGroup(&Group{
				Name: group,
				Arn:  "arn:aws:iam::123456789012:group/" + group,
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

func RootModel(m tea.Model) Groups {
	return Groups{
		model: m,
	}
}
