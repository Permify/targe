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
		var model tea.Model

		if args[0] == USERS.String() {
			if len(args) > 1 {
				if len(args) > 2 {
					if len(args) > 3 {
						// TODO: is it arn or name?
						result := Result(User{
							Name: args[1],
							Arn:  "arn:aws:iam::123456789012:user/" + args[1],
						}, Policy{
							Name: args[3],
							Arn:  "arn:aws:iam::aws:policy/" + args[3],
						}, args[2])

						model = &result
					} else {
						// TODO: is it arn or name?
						policies := Policies(User{
							Name: args[1],
							Arn:  "arn:aws:iam::123456789012:user/" + args[1],
						}, args[2])
						model = &policies
					}
				} else {
					// TODO: is it arn or name?
					actions := Actions(User{
						Name: args[1],
						Arn:  "arn:aws:iam::123456789012:user/" + args[1],
					})
					model = &actions
				}
			} else {
				users := Users()
				model = &users
			}
		} else {
			users := Users()
			model = &users
		}

		p := tea.NewProgram(RootModel(model), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		return nil
	}
}

func RootModel(m tea.Model) Aws {
	return Aws{
		model: m,
	}
}

func Switch(model tea.Model, width, height int) (tea.Model, tea.Cmd) {
	if width == 0 && height == 0 {
		return model, model.Init()
	}

	return model.Update(tea.WindowSizeMsg{
		Width:  width,
		Height: height,
	})
}
