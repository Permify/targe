package users

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	state *State
	list  list.Model
}

func CustomPolicyOptions(state *State) CustomPolicyOptionModel {
	items := []list.Item{
		CustomPolicyOption{Name: "Without Resource (without_resource)", Desc: "Applies globally without a resource."},
		CustomPolicyOption{Name: "With Resource (with_resource)", Desc: "Scoped to a specific resource."},
	}
	var m CustomPolicyOptionModel
	m.state = state
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
			m.state.SetPolicyOption(&option)
			return Switch(m.state.FindFlow(), 0, 0)
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
