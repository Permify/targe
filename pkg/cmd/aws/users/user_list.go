package users

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var usersStyle = lipgloss.NewStyle().Margin(1, 2)

type User struct {
	Arn  string
	Name string
}

func (i User) Title() string       { return i.Name }
func (i User) Description() string { return i.Arn }
func (i User) FilterValue() string { return i.Arn }

type UserListModel struct {
	state *State
	list  list.Model
}

func UserList(state *State) UserListModel {
	var items []list.Item

	users := []User{
		{
			Name: "Alice",
			Arn:  "arn:aws:iam::123456789012:user/Alice",
		},
		{
			Name: "Bob",
			Arn:  "arn:aws:iam::123456789012:user/Bob",
		},
	}

	for _, user := range users {
		items = append(items, User{
			Name: user.Name,
			Arn:  user.Arn,
		})
	}

	var m UserListModel
	m.state = state
	m.list.Title = "Users"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m UserListModel) Init() tea.Cmd {
	return nil
}

func (m UserListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			user := m.list.SelectedItem().(User)
			m.state.SetUser(&user)
			return Switch(m.state.Next(), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := usersStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m UserListModel) View() string {
	return usersStyle.Render(m.list.View())
}
