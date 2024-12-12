package aws

import (
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/config"
	"github.com/Permify/kivo/pkg/cmd/aws/groups"
	"github.com/Permify/kivo/pkg/cmd/aws/roles"
	"github.com/Permify/kivo/pkg/cmd/aws/users"
)

// NewAwsCommand -
func NewAwsCommand(cfg *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "aws",
		Short: "",
		Long:  ``,
	}

	command.AddCommand(users.NewUsersCommand(cfg))
	command.AddCommand(roles.NewRolesCommand(cfg))
	command.AddCommand(groups.NewGroupsCommand(cfg))

	return command
}
