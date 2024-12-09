package aws

import (
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/pkg/cmd/aws/users"
)

// NewAwsCommand -
func NewAwsCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "aws",
		Short: "",
		Long:  ``,
	}

	command.AddCommand(users.NewUsersCommand())

	return command
}
