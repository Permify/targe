package users

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Permify/kivo/internal/aws"
)

var resourcesStyle = lipgloss.NewStyle().Margin(1, 2)

type Resource struct {
	Arn  string
	Name string
}

func (i Resource) Title() string       { return i.Name }
func (i Resource) Description() string { return i.Arn }
func (i Resource) FilterValue() string { return i.Arn }

type ResourceListModel struct {
	state *State
	list  list.Model
	err   error
}

func ResourceList(state *State) ResourceListModel {
	var items []list.Item

	resources, err := aws.ListResources(state.GetService().Name)

	for _, resource := range resources {
		items = append(items, Resource{
			Name: resource.Name,
			Arn:  resource.Arn,
		})
	}

	var m ResourceListModel
	m.err = err
	m.state = state
	m.list.Title = fmt.Sprintf("Resources for %s", state.GetService().Title())
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m ResourceListModel) Init() tea.Cmd {
	return nil
}

func (m ResourceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.list.Items()) != 0 {
				resource := m.list.SelectedItem().(Resource)
				m.state.SetResource(&resource)
			} else {
				m.state.SetService(nil)
				m.state.SetPolicyOption(nil)
			}
			return Switch(m.state.Next(), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := policiesStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ResourceListModel) View() string {
	if m.err != nil {
		return resourcesStyle.Render(m.err.Error())
	}

	if len(m.list.Items()) == 0 {
		return resourcesStyle.Render("No resources found.")
	}

	// Style for the resource list
	return resourcesStyle.Render(m.list.View())
}
