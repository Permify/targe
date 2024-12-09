package users

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var userActionsStyle = lipgloss.NewStyle().Margin(1, 2)

type Action struct {
	Id   string
	Name string
	Desc string
}

func (i Action) Title() string       { return i.Name }
func (i Action) Description() string { return i.Desc }
func (i Action) FilterValue() string { return i.Name }

type ActionListModel struct {
	state *State
	list  list.Model
}

func ActionList(state *State) ActionListModel {
	items := []list.Item{
		Action{Id: AttachPolicySlug, Name: ReachableActions[AttachPolicySlug].Name, Desc: ReachableActions[AttachPolicySlug].Desc},
		Action{Id: DetachPolicySlug, Name: ReachableActions[DetachPolicySlug].Name, Desc: ReachableActions[DetachPolicySlug].Desc},
		Action{Id: AddToGroupSlug, Name: ReachableActions[AddToGroupSlug].Name, Desc: ReachableActions[AddToGroupSlug].Desc},
		Action{Id: RemoveFromGroupSlug, Name: ReachableActions[RemoveFromGroupSlug].Name, Desc: ReachableActions[RemoveFromGroupSlug].Desc},
		Action{Id: AttachCustomPolicySlug, Name: ReachableActions[AttachCustomPolicySlug].Name, Desc: ReachableActions[AttachCustomPolicySlug].Desc},
	}
	var m ActionListModel
	m.state = state
	m.list.Title = "Actions"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m ActionListModel) Init() tea.Cmd {
	return nil
}

func (m ActionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			action := m.list.SelectedItem().(Action)
			m.state.SetAction(&action)
			return Switch(m.state.FindFlow(), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := userActionsStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ActionListModel) View() string {
	return userActionsStyle.Render(m.list.View())
}
