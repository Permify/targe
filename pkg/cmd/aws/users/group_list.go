package users

import (
	"context"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	internalaws "github.com/Permify/kivo/internal/aws"
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
	err   error
}

func GroupList(state *State) GroupListModel {
	var items []list.Item
	var m GroupListModel

	groups, err := internalaws.ListGroups(context.Background(), state.awsConfig)

	userGroups, err := internalaws.ListGroupsForUser(context.Background(), state.awsConfig, state.user.Name)
	m.err = err

	switch state.operation.Id {
	case AddToGroupSlug:
		for _, group := range groups.Groups {
			if !slices.Contains(userGroups, aws.ToString(group.GroupName)) {
				items = append(items, Group{
					Name: aws.ToString(group.GroupName),
					Arn:  aws.ToString(group.Arn),
				})
			}
		}
	case RemoveFromGroupSlug:
		for _, group := range groups.Groups {
			if slices.Contains(userGroups, aws.ToString(group.GroupName)) {
				items = append(items, Group{
					Name: aws.ToString(group.GroupName),
					Arn:  aws.ToString(group.Arn),
				})
			}
		}
	}

	m.err = err
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
