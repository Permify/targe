package groups

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	internalaws "github.com/Permify/kivo/internal/aws"
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

		// Load the AWS configuration
		awscfg, err := awsconfig.LoadDefaultConfig(context.Background())
		if err != nil {
			return err
		}

		api := internalaws.NewApi(awscfg)
		state := &State{}

		if group != "" {
			awsgroup, err := api.FindGroup(context.Background(), group)
			if err != nil {
				return err
			}

			state.SetGroup(&Group{
				Name: aws.ToString(awsgroup.Group.GroupName),
				Arn:  aws.ToString(awsgroup.Group.Arn),
			})
		}

		if operation != "" {
			// Check if the operation exists in the ReachableOperations map
			op, exists := ReachableOperations[operation]
			if !exists {
				return fmt.Errorf("Operation '%s' does not exist in ReachableOperations\n", operation)
			}

			state.SetOperation(&Operation{
				Id:   op.Id,
				Name: op.Name,
				Desc: op.Desc,
			})
		}

		if policy != "" {
			awspolicy, err := api.FindPolicy(context.Background(), policy)
			if err != nil {
				return err
			}

			state.SetPolicy(&Policy{
				Name: aws.ToString(awspolicy.Policy.PolicyName),
				Arn:  aws.ToString(awspolicy.Policy.Arn),
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
			// Check if the operation exists in the ReachableOperations map
			op, exists := ReachableCustomPolicyOptions[policyOption]
			if !exists {
				return fmt.Errorf("Policy options '%s' does not exist in ReachableCustomPolicyOptions\n", policyOption)
			}

			state.SetPolicyOption(&CustomPolicyOption{
				Name: op.Name,
				Desc: op.Desc,
			})
		}

		p := tea.NewProgram(RootModel(state.Next(api)), tea.WithAltScreen())
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
