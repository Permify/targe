package aws

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/pkg/cmd/aws/users"
)

type Root string

const (
	USERS  Root = "users"
	GROUPS Root = "groups"
	ROLES  Root = "roles"
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
		var state users.State

		if args[0] == USERS.String() {
			if len(args) > 1 {
				if len(args) > 2 {
					if len(args) > 3 {
						// TODO: is it arn or name?

						state.SetUser(&users.User{
							Name: args[1],
							Arn:  "arn:aws:iam::123456789012:user/" + args[1],
						})

						state.SetPolicy(&users.Policy{
							Name: args[3],
							Arn:  "arn:aws:iam::aws:policy/" + args[3],
						})

					} else {
						// TODO: is it arn or name?

						state.SetUser(&users.User{
							Name: args[1],
							Arn:  "arn:aws:iam::123456789012:user/" + args[1],
						})

						state.SetAction(&users.UserAction{
							Name: args[2],
							Desc: users.ReachableActions[args[2]].Desc,
						})
					}
				} else {
					// TODO: is it arn or name?
					state.SetUser(&users.User{
						Name: args[1],
						Arn:  "arn:aws:iam::123456789012:user/" + args[1],
					})
				}
			}
		}

		p := tea.NewProgram(RootModel(state.FindFlow()), tea.WithAltScreen())
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
