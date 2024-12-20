package groups

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Permify/kivo/pkg/aws/models"
)

type OperationList struct {
	controller *Controller
	spinner    spinner.Model
	loading    bool
	list       list.Model
	err        error
}

func NewOperationList(controller *Controller) OperationList {
	sp := spinner.New()
	sp.Style = spinnerStyle
	sp.Spinner = spinner.Pulse

	view := OperationList{
		controller: controller,
		spinner:    sp,
		loading:    true,
	}

	view.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	view.list.Title = "Operations"
	return view
}

func (m OperationList) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.controller.LoadOperations())
}

func (m OperationList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.loading {
				option := m.list.SelectedItem().(models.Operation)
				m.controller.State.SetOperation(&option)
				return Switch(m.controller.Next(), m.list.Width(), m.list.Height())
			}
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case OperationLoadedMsg:
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

func (m OperationList) View() string {
	if m.err != nil {
		return listStyle.Render(m.err.Error())
	}

	if m.loading {
		return listStyle.Render(m.spinner.View() + " Loading...")
	}

	return listStyle.Render(m.list.View())
}
