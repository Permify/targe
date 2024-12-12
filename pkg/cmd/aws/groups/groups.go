package groups

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

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

	return command
}

func groups(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
