package aws

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var resourcesStyle = lipgloss.NewStyle().Margin(1, 2)

type Resource struct {
	Arn  string
	Name string
}

func (i Resource) Title() string       { return i.Name }
func (i Resource) Description() string { return i.Arn }
func (i Resource) FilterValue() string { return i.Arn }

type ResourcesModel struct {
	service string
	user    User
	action  string
	list    list.Model
}

func Resources(user User, action, service string) ResourcesModel {
	var items []list.Item
	resources := []Resource{
		{
			Name: "S3 Bucket",
			Arn:  "arn:aws:s3:::my_bucket",
		},
		{
			Name: "DynamoDB Table",
			Arn:  "arn:aws:dynamodb:us-west-2:123456789012:table/my_table",
		},
	}

	for _, resource := range resources {
		items = append(items, Resource{
			Name: resource.Name,
			Arn:  resource.Arn,
		})
	}

	var m ResourcesModel
	m.user = user
	m.action = action
	m.service = service
	m.list.Title = fmt.Sprintf("Resources for %s", service)
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m ResourcesModel) Init() tea.Cmd {
	return nil
}

func (m ResourcesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			resource := m.list.SelectedItem().(Resource)
			createPolicy := CreatePolicy(m.user, resource)
			return Switch(&createPolicy, m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := policiesStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ResourcesModel) View() string {
	return resourcesStyle.Render(m.list.View())
}
