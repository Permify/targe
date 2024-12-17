package users

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	internalaws "github.com/Permify/kivo/internal/aws"
	"github.com/Permify/kivo/internal/requirements/aws"
)

var servicesStyle = lipgloss.NewStyle().Margin(1, 2)

type Service struct {
	Name string
	Desc string
}

func (i Service) Title() string       { return i.Name }
func (i Service) Description() string { return i.Desc }
func (i Service) FilterValue() string { return i.Name }

type ServiceListModel struct {
	api   *internalaws.Api
	state *State
	list  list.Model
	err   error
}

func ServiceList(api *internalaws.Api, state *State) ServiceListModel {
	var m ServiceListModel
	m.api = api
	m.state = state

	t := aws.Types{}
	services, err := t.GetServices()

	var items []list.Item

	for _, service := range services {
		items = append(items, Service{
			Name: service.Name,
			Desc: service.Description,
		})
	}

	m.err = err
	m.list.Title = "Services"
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m ServiceListModel) Init() tea.Cmd {
	return nil
}

func (m ServiceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.list.Items()) != 0 {
				service := m.list.SelectedItem().(Service)
				m.state.SetService(&service)
			}
			return Switch(m.state.Next(m.api), m.list.Width(), m.list.Height())
		}
	case tea.WindowSizeMsg:
		h, v := servicesStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ServiceListModel) View() string {
	if m.err != nil {
		return servicesStyle.Render(m.err.Error())
	}

	if len(m.list.Items()) == 0 {
		return servicesStyle.Render("No services found.")
	}

	return servicesStyle.Render(m.list.View())
}
