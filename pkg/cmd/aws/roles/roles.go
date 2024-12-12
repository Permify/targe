package roles

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

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

	return command
}

func roles(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
