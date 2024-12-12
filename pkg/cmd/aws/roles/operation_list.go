package roles

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var userOperationsStyle = lipgloss.NewStyle().Margin(1, 2)

type Operation struct {
	Id   string
	Name string
	Desc string
}

func (i Operation) Title() string       { return i.Name }
func (i Operation) Description() string { return i.Desc }
func (i Operation) FilterValue() string { return i.Name }

type OperationListModel struct {
	state *State
	list  list.Model
}

func OperationList(state *State) OperationListModel {
	items := []list.Item{
		Operation{Id: AttachPolicySlug, Name: ReachableOperations[AttachPolicySlug].Name, Desc: ReachableOperations[AttachPolicySlug].Desc},
		Operation{Id: DetachPolicySlug, Name: ReachableOperations[DetachPolicySlug].Name, Desc: ReachableOperations[DetachPolicySlug].Desc},
		Operation{Id: AttachCustomPolicySlug, Name: ReachableOperations[AttachCustomPolicySlug].Name, Desc: ReachableOperations[AttachCustomPolicySlug].Desc},
	}
	var m OperationListModel
	m.state = state
	m.list.Title = "Operations"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m OperationListModel) Init() tea.Cmd {
	return nil
}

func (m OperationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			option := m.list.SelectedItem().(Operation)
			m.state.SetOperation(&option)
			return Switch(m.state.Next(), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := userOperationsStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m OperationListModel) View() string {
	return userOperationsStyle.Render(m.list.View())
}
