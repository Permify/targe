package users

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var customPolicyOptionsStyle = lipgloss.NewStyle().Margin(1, 2)

type CustomPolicyOption struct {
	Id   string
	Name string
	Desc string
}

func (i CustomPolicyOption) Title() string       { return i.Name }
func (i CustomPolicyOption) Description() string { return i.Desc }
func (i CustomPolicyOption) FilterValue() string { return i.Name }

type CustomPolicyOptionListModel struct {
	state *State
	list  list.Model
}

func CustomPolicyOptionList(state *State) CustomPolicyOptionListModel {
	items := []list.Item{
		CustomPolicyOption{Id: WithoutResourceSlug, Name: ReachableCustomPolicyOptions[WithoutResourceSlug].Name, Desc: ReachableCustomPolicyOptions[WithoutResourceSlug].Desc},
		CustomPolicyOption{Id: WithResourceSlug, Name: ReachableCustomPolicyOptions[WithResourceSlug].Name, Desc: ReachableCustomPolicyOptions[WithResourceSlug].Desc},
	}
	var m CustomPolicyOptionListModel
	m.state = state
	m.list.Title = "Custom Policy Options"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m CustomPolicyOptionListModel) Init() tea.Cmd {
	return nil
}

func (m CustomPolicyOptionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			option := m.list.SelectedItem().(CustomPolicyOption)
			m.state.SetPolicyOption(&option)
			return Switch(m.state.Next(), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := userActionsStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CustomPolicyOptionListModel) View() string {
	return customPolicyOptionsStyle.Render(m.list.View())
}
