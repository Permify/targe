package roles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Permify/kivo/pkg/aws/models"
)

type RoleList struct {
	controller *Controller
	spinner    spinner.Model
	loading    bool
	list       list.Model
	err        error
}

func NewRoleList(controller *Controller) RoleList {
	sp := spinner.New()
	sp.Style = spinnerStyle
	sp.Spinner = spinner.Pulse

	view := RoleList{
		controller: controller,
		spinner:    sp,
		loading:    true,
	}

	view.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	view.list.Title = "Roles"
	return view
}

func (m RoleList) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.controller.LoadRoles())
}

func (m RoleList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.loading {
				role := m.list.SelectedItem().(models.Role)
				m.controller.State.SetRole(&role)
				return Switch(m.controller.Next(), m.list.Width(), m.list.Height())
			}
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case RoleLoadedMsg:
		// Update list with loaded users
		m.loading = false
		m.list.SetItems(msg.List)
	case FailedMsg:
		// Handle error
		m.loading = false
		m.err = msg.Err
	}

	// Update spinner if loading
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Update list if not loading
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m RoleList) View() string {
	if m.err != nil {
		return listStyle.Render(m.err.Error())
	}

	if m.loading {
		return listStyle.Render(m.spinner.View() + " Loading...")
	}

	return listStyle.Render(m.list.View())
}
