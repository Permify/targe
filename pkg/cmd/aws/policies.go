package aws

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var policiesStyle = lipgloss.NewStyle().Margin(1, 2)

type Policy struct {
	Arn  string
	Name string
}

func (i Policy) Title() string       { return i.Name }
func (i Policy) Description() string { return i.Arn }
func (i Policy) FilterValue() string { return i.Arn }

type PoliciesModel struct {
	list list.Model
}

func Policies() PoliciesModel {
	var items []list.Item
	policies := []Policy{
		{
			Name: "AdministratorAccess",
			Arn:  "arn:aws:iam::aws:policy/AdministratorAccess",
		},
		{
			Name: "PowerUserAccess",
			Arn:  "arn:aws:iam::aws:policy/PowerUserAccess",
		},
		{
			Name: "ReadOnlyAccess",
			Arn:  "arn:aws:iam::aws:policy/ReadOnlyAccess",
		},
		{
			Name: "SecurityAudit",
			Arn:  "arn:aws:iam::aws:policy/SecurityAudit",
		},
		{
			Name: "NetworkAdministrator",
			Arn:  "arn:aws:iam::aws:policy/NetworkAdministrator",
		},
	}

	for _, policy := range policies {
		items = append(items, User{
			Name: policy.Name,
			Arn:  policy.Arn,
		})
	}

	var m PoliciesModel
	m.list.Title = "Policies"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m PoliciesModel) Init() tea.Cmd {
	return nil
}

func (m PoliciesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := policiesStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PoliciesModel) View() string {
	return usersStyle.Render(m.list.View())
}
