package users

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Permify/kivo/internal/aws"
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
	err   error
}

func UserList(state *State) UserListModel {
	var items []list.Item

	output, err := aws.ListUsers(context.Background(), state.awsConfig)

	for _, user := range output.Users {
		items = append(items, User{
			Name: *user.UserName,
			Arn:  *user.Arn,
		})
	}

	var m UserListModel
	m.err = err
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
	if m.err != nil {
		return usersStyle.Render(m.err.Error())
	}

	return usersStyle.Render(m.list.View())
}
