package roles

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

type Roles struct {
	model tea.Model
}

func (m Roles) Init() tea.Cmd {
	return m.model.Init() // rest methods are just wrappers for the model's methods
}

func (m Roles) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m Roles) View() string {
	return m.model.View()
}

// NewRolesCommand -
func NewRolesCommand(cfg *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "roles",
		Short: "",
		RunE:  roles(cfg),
	}

	f := command.Flags()

	f.String("role", "", "role")
	f.String("operation", "", "operation")
	f.String("policy", "", "policy")
	f.String("resource", "", "resource")
	f.String("service", "", "service")
	f.String("policy-option", "", "policy option")

	// SilenceUsage is set to true to suppress usage when an error occurs
	command.SilenceUsage = true

	command.PreRun = func(cmd *cobra.Command, args []string) {
		RegisterRolesFlags(f)
	}

	return command
}

func roles(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// get min coverage from viper
		role := viper.GetString("role")
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

		if role != "" {
			awsrole, err := api.FindRole(context.Background(), role)
			if err != nil {
				return err
			}

			state.SetRole(&Role{
				Name: aws.ToString(awsrole.Role.RoleName),
				Arn:  aws.ToString(awsrole.Role.Arn),
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

func RootModel(m tea.Model) Roles {
	return Roles{
		model: m,
	}
}
