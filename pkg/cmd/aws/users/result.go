package users

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss/table"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	internalaws "github.com/Permify/kivo/internal/aws"
)

const maxWidth = 100

var (
	red   = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	blue  = lipgloss.Color("212")
	green = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(blue).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type ResultModel struct {
	api    *internalaws.Api
	state  *State
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
}

func Result(api *internalaws.Api, state *State) ResultModel {
	m := ResultModel{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.api = api
	m.state = state

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

func (m ResultModel) Init() tea.Cmd {
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m ResultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ResultModel) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render("DONE") + "\n\n"
	default:

		var rows [][]string

		if m.state.user != nil {
			rows = append(rows, []string{"User", m.state.user.Name, m.state.user.Arn})
		}

		if m.state.operation != nil {
			rows = append(rows, []string{"Operation", m.state.operation.Name, m.state.operation.Desc})
		}

		if m.state.group != nil {
			rows = append(rows, []string{"Group", m.state.group.Name, m.state.group.Arn})
		}

		if m.state.service != nil {
			rows = append(rows, []string{"Service", m.state.service.Name, m.state.service.Desc})
		}

		if m.state.resource != nil {
			rows = append(rows, []string{"Resource", m.state.resource.Name, m.state.resource.Arn})
		}

		if m.state.policy != nil {
			if len(m.state.policy.Document) > 0 {
				// Marshal with indent
				indentedJSON, err := json.MarshalIndent(m.state.policy.Document, "", "  ")
				if err != nil {
					fmt.Println("Error:", err)
				}

				rows = append(rows, []string{"Policy", m.state.policy.Name, string(indentedJSON)})
			} else {
				rows = append(rows, []string{"Policy", m.state.policy.Name, m.state.policy.Arn})
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

func (m ResultModel) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m ResultModel) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(blue),
	)
}

func (m ResultModel) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}
