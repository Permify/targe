package groups

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var groupsStyle = lipgloss.NewStyle().Margin(1, 2)

type Group struct {
	Arn  string
	Name string
}

func (i Group) Title() string       { return i.Name }
func (i Group) Description() string { return i.Arn }
func (i Group) FilterValue() string { return i.Arn }

type GroupListModel struct {
	state *State
	list  list.Model
}

func GroupList(state *State) GroupListModel {
	var items []list.Item

	groups := []Group{
		{
			Name: "Group 1",
			Arn:  "arn:aws:iam::123456789012:group/Group1",
		},
		{
			Name: "Group 2",
			Arn:  "arn:aws:iam::123456789012:group/Group2",
		},
	}

	for _, group := range groups {
		items = append(items, Group{
			Name: group.Name,
			Arn:  group.Arn,
		})
	}

	var m GroupListModel
	m.state = state
	m.list.Title = "Groups"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m GroupListModel) Init() tea.Cmd {
	return nil
}

func (m GroupListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			group := m.list.SelectedItem().(Group)
			m.state.SetGroup(&group)
			return Switch(m.state.Next(), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := groupsStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m GroupListModel) View() string {
	return groupsStyle.Render(m.list.View())
}
