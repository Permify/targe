package aws

import (
	`github.com/charmbracelet/bubbles/list`
	tea `github.com/charmbracelet/bubbletea`
	`github.com/charmbracelet/lipgloss`
)

var customPolicyOptionsStyle = lipgloss.NewStyle().Margin(1, 2)

type CustomPolicyOption struct {
	Name string
	Desc string
}

func (i CustomPolicyOption) Title() string       { return i.Name }
func (i CustomPolicyOption) Description() string { return i.Desc }
func (i CustomPolicyOption) FilterValue() string { return i.Name }

type CustomPolicyOptionModel struct {
	user User
	list list.Model
}

func CustomPolicyOptions(user User) CustomPolicyOptionModel {
	items := []list.Item{
		CustomPolicyOption{Name: "Policy Without Resource (policy_without_resource)", Desc: "Applies globally without a resource."},
		CustomPolicyOption{Name: "Policy With Resource (policy_with_resource)", Desc: "Scoped to a specific resource."},
	}
	var m CustomPolicyOptionModel
	m.user = user
	m.list.Title = "Custom Policy Options"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m CustomPolicyOptionModel) Init() tea.Cmd {
	return nil
}
func (m CustomPolicyOptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			option := m.list.SelectedItem().(CustomPolicyOption)

			if option.Name == "Policy Without Resource (policy_without_resource)" {
				createPolicy := CreatePolicy(m.user, Resource{})
				return Switch(&createPolicy, 0, 0)
			}

			if option.Name == "Policy With Resource (policy_with_resource)" {
				createPolicy := Services(m.user, "custom_with_resource")
				return Switch(&createPolicy, m.list.Width(), m.list.Height())
			}
		}
	case tea.WindowSizeMsg:
		h, v := userActionsStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CustomPolicyOptionModel) View() string {
	return customPolicyOptionsStyle.Render(m.list.View())
}
