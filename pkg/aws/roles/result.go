package roles

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Result struct {
	controller *Controller
	lg         *lipgloss.Renderer
	styles     *Styles
	form       *huh.Form
	width      int
}

func NewResult(controller *Controller) Result {
	m := Result{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.controller = controller

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					return nil
				}).
				Affirmative("Yes").
				Negative("No"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m Result) Init() tea.Cmd {
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m Result) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Result) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render("DONE") + "\n\n"
	default:

		var rows [][]string

		if m.controller.State.role != nil {
			rows = append(rows, []string{"Role", m.controller.State.role.Name, m.controller.State.role.Arn})
		}

		if m.controller.State.operation != nil {
			rows = append(rows, []string{"Operation", m.controller.State.operation.Name, m.controller.State.operation.Desc})
		}

		if m.controller.State.service != nil {
			rows = append(rows, []string{"Service", m.controller.State.service.Name, m.controller.State.service.Desc})
		}

		if m.controller.State.resource != nil {
			rows = append(rows, []string{"Resource", m.controller.State.resource.Name, m.controller.State.resource.Arn})
		}

		if m.controller.State.policy != nil {
			if len(m.controller.State.policy.Document) > 0 {
				// Marshal with indent
				indentedJSON, err := json.MarshalIndent(m.controller.State.policy.Document, "", "  ")
				if err != nil {
					fmt.Println("Error:", err)
				}

				rows = append(rows, []string{"Policy", m.controller.State.policy.Name, string(indentedJSON)})
			} else {
				rows = append(rows, []string{"Policy", m.controller.State.policy.Name, m.controller.State.policy.Arn})
			}
		}

		t := table.New().
			Border(lipgloss.HiddenBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				if col == 0 {
					return m.styles.Base.
						Foreground(blue).
						Bold(true)
				}
				return m.styles.Base
			}).
			Rows(rows...)

		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		errors := m.form.Errors()
		header := m.appBoundaryView("Overview")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinVertical(lipgloss.Top, t.Render(), form)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Result) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Result) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(blue),
	)
}

func (m Result) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}
