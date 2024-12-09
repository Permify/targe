package users

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var resultStyle = lipgloss.NewStyle().Margin(1, 2)

type ResultModel struct {
	state *State
}

func Result(state *State) ResultModel {
	return ResultModel{
		state: state,
	}
}

func (m ResultModel) Init() tea.Cmd {
	return nil
}

func (m ResultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // Quit
			return m, tea.Quit
		case "enter": // Quit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ResultModel) View() string {
	return resultStyle.Render(fmt.Sprintf(
		`Are you sure you want to proceed?
			"User:   %s
			"Action: %s
			"Policy: %s
			"Press [Y] to confirm or [N] to cancel.`,
		m.state.GetUser().Name, m.state.GetAction().Name, m.state.GetPolicy().Name,
	))
}
