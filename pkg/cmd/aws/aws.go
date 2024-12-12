package aws

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/config"
	"github.com/Permify/kivo/pkg/cmd/aws/groups"
	"github.com/Permify/kivo/pkg/cmd/aws/roles"
	"github.com/Permify/kivo/pkg/cmd/aws/users"
	"github.com/Permify/kivo/pkg/cmd/common"
)

// NewAwsCommand -
func NewAwsCommand(cfg *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "aws",
		Short: "",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
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
		},
	}

	command.AddCommand(users.NewUsersCommand(cfg))
	command.AddCommand(roles.NewRolesCommand(cfg))
	command.AddCommand(groups.NewGroupsCommand(cfg))

	return command
}

func folderExists(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
