package aws

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type Root string

const (
	USERS    Root = "users"
	POLICIES Root = "policies"
)

func (r Root) String() string {
	return string(r)
}

// NewAwsCommand -
func NewAwsCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "aws",
		Short: "",
		RunE:  aws(),
	}
	return command
}

type Aws struct {
	model tea.Model
}

func (m Aws) Init() tea.Cmd {
	return m.model.Init() // rest methods are just wrappers for the model's methods
}

func (m Aws) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m Aws) View() string {
	return m.model.View()
}

func aws() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(RootScreen(Root(args[0])), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		return nil
	}
}

func RootScreen(m Root) Aws {
	var root tea.Model

	switch m {
	case USERS:
		users := Users()
		root = &users
	default:
		users := Users()
		root = &users
	}

	return Aws{
		model: root,
	}
}

func Switch(model tea.Model, width, height int) (tea.Model, tea.Cmd) {
	return model.Update(tea.WindowSizeMsg{
		Width:  width,
		Height: height,
	})
}
