package aws

import (
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/config"
)

// NewAwsCommand -
func NewAwsCommand(cfg *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "aws",
		Short: "",
		Long:  ``,
	}

	command.AddCommand(NewUsersCommand(cfg))
	command.AddCommand(NewRolesCommand(cfg))
	command.AddCommand(NewGroupsCommand(cfg))

	return command
}
