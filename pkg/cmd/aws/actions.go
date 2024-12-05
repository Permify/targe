package aws

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var actionsStyle = lipgloss.NewStyle().Margin(1, 2)

type Action struct {
	Name string
	Desc string
}

func (i Action) Title() string       { return i.Name }
func (i Action) Description() string { return i.Desc }
func (i Action) FilterValue() string { return i.Name }

type ActionsModel struct {
	user        User
	list        list.Model
	sharedState SharedState
}

func Actions(user User) ActionsModel {
	items := []list.Item{
		Action{Name: "Attach Policy", Desc: "Assign a policy to the user."},
		Action{Name: "Detach Policy", Desc: "Remove a policy from the user."},
		Action{Name: "Add to Group", Desc: "Include the user in a group."},
		Action{Name: "Remove from Group", Desc: "Exclude the user from a group."},
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
		}
	case tea.WindowSizeMsg:
		h, v := actionsStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ActionsModel) View() string {
	return actionsStyle.Render(m.list.View())
}
