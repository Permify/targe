package aws

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var userActionsStyle = lipgloss.NewStyle().Margin(1, 2)

type UserAction struct {
	Name string
	Desc string
}

func (i UserAction) Title() string       { return i.Name }
func (i UserAction) Description() string { return i.Desc }
func (i UserAction) FilterValue() string { return i.Name }

type ActionsModel struct {
	user User
	list list.Model
}

func Actions(user User) ActionsModel {
	items := []list.Item{
		UserAction{Name: "Attach Policy (attach_policy)", Desc: "Assign a policy to the user."},
		UserAction{Name: "Detach Policy (detach_policy)", Desc: "Remove a policy from the user."},
		UserAction{Name: "Add to Group (add_to_group)", Desc: "Include the user in a group."},
		UserAction{Name: "Remove from Group (remove_from_group)", Desc: "Exclude the user from a group."},
		UserAction{Name: "Custom Policy (custom_policy)", Desc: "Crete and attach custom policy."},
	}
	var m ActionsModel
	m.user = user
	m.list.Title = "Actions"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m ActionsModel) Init() tea.Cmd {
	return nil
}

func (m ActionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			action := m.list.SelectedItem().(UserAction)

			if action.Name == "Custom Policy (custom_policy)" {
				customPolicy := CustomPolicyOptions(m.user)
				return Switch(&customPolicy, m.list.Width(), m.list.Height())
			}

			policies := Policies(m.user, action.Name)
			return Switch(&policies, m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := userActionsStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ActionsModel) View() string {
	return userActionsStyle.Render(m.list.View())
}
