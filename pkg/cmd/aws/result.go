package aws

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var resultStyle = lipgloss.NewStyle().Margin(1, 2)

type ResultModel struct {
	User   User
	Policy Policy
	Action string
}

func Result(user User, policy Policy, action string) ResultModel {
	return ResultModel{
		User:   user,
		Policy: policy,
		Action: action,
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
		"Are you sure you want to proceed?\n\n"+
			"User:   %s\n"+
			"Action: %s\n"+
			"Policy: %s\n\n"+
			"Press [Y] to confirm or [N] to cancel.",
		m.User.Name, m.Action, m.Policy.Name,
	))
}
