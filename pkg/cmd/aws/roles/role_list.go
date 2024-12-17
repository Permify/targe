package roles

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	internalaws "github.com/Permify/kivo/internal/aws"
)

var rolesStyle = lipgloss.NewStyle().Margin(1, 2)

type Role struct {
	Arn  string
	Name string
}

func (i Role) Title() string       { return i.Name }
func (i Role) Description() string { return i.Arn }
func (i Role) FilterValue() string { return i.Arn }

type RoleListModel struct {
	api   *internalaws.Api
	state *State
	list  list.Model
	err   error
}

func RoleList(api *internalaws.Api, state *State) RoleListModel {
	var items []list.Item

	output, err := api.ListRoles(context.Background())

	for _, role := range output.Roles {
		items = append(items, Role{
			Name: *role.RoleName,
			Arn:  *role.Arn,
		})
	}

	var m RoleListModel
	m.err = err
	m.api = api
	m.state = state
	m.list.Title = "Roles"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m RoleListModel) Init() tea.Cmd {
	return nil
}

func (m RoleListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			role := m.list.SelectedItem().(Role)
			m.state.SetRole(&role)
			return Switch(m.state.Next(m.api), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := rolesStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m RoleListModel) View() string {
	return rolesStyle.Render(m.list.View())
}
