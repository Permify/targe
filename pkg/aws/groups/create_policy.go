package groups

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Permify/kivo/internal/ai"
	"github.com/Permify/kivo/pkg/aws/models"
)

type CreatePolicy struct {
	controller  *Controller
	lg          *lipgloss.Renderer
	styles      *Styles
	form        *huh.Form
	senderStyle lipgloss.Style
	err         error
	width       int
	message     *string
	done        *bool
	result      string
}

func NewCreatePolicy(controller *Controller) CreatePolicy {
	m := CreatePolicy{controller: controller, width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	doneInitialValue := false
	m.done = &doneInitialValue

	messageInitialValue := ""
	m.message = &messageInitialValue

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewText().Key("message").
				Title("Describe Your Policy").
				Value(m.message),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Value(m.done).
				Affirmative("Yes").
				Negative("Refresh"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)

	return m
}

func (m CreatePolicy) Init() tea.Cmd {
	return m.form.Init()
}

func (m CreatePolicy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()

	case tea.KeyMsg:

		if msg.String() == "esc" || msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

		// Check if the "Refresh" or "Done" button was selected
		if msg.String() == "enter" {
			if m.done != nil && *m.done {
				return Switch(m.controller.Next(), 0, 0)
			} else {
				var resourceArn *string = nil
				if m.controller.State.GetService() != nil {
					resourceArn = &m.controller.State.GetResource().Arn
				}

				if m.message == nil {
					m.err = errors.New("Please provide a message")
				}

				policy, err := ai.GeneratePolicy(m.controller.openAiApiKey, *m.message, resourceArn)
				if err != nil {
					m.err = err
				}

				policyJson, err := json.MarshalIndent(policy, "", "\t")
				if err != nil {
					m.err = err
				}

				m.result = string(policyJson)

				m.controller.State.SetPolicy(&models.Policy{
					Arn:      "new",
					Name:     policy.Id,
					Document: string(policyJson),
				})

				m.reinitializeForm()
			}
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m CreatePolicy) View() string {
	s := m.styles

	v := strings.TrimSuffix(m.form.View(), "\n\n")
	form := m.lg.NewStyle().Margin(1, 0).Render(v)

	var titles string
	if m.controller.State.GetService() != nil && m.controller.State.GetResource() != nil {
		titles = lipgloss.JoinVertical(lipgloss.Left,
			s.ServiceNameHeader.Render("Service Name: "+m.controller.State.GetService().Name),
			s.ResourceArnHeader.Render("Resource ARN: "+m.controller.State.GetResource().Arn),
		)
	}

	// Status (right side)
	var status string
	{
		buildInfo := "(None)"

		if m.result != "" {
			buildInfo = m.result
		}

		const statusWidth = 60
		statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
		status = s.Status.
			Height(lipgloss.Height(form)).
			Width(statusWidth).
			MarginLeft(statusMarginLeft).
			Render(s.StatusHeader.Render("Policy") + "\n" +
				buildInfo)
	}

	errors := m.form.Errors()
	header := lipgloss.JoinVertical(lipgloss.Top,
		m.appBoundaryView("Custom Policy Generator"),
		titles,
	)
	if len(errors) > 0 {
		header = m.appErrorBoundaryView(m.errorView())
	}
	body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

	footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
	if len(errors) > 0 {
		footer = m.appErrorBoundaryView("")
	}

	return s.Base.Render(header + "\n" + body + "\n\n" + footer)
}

func (m CreatePolicy) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m CreatePolicy) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m CreatePolicy) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func (m *CreatePolicy) reinitializeForm() {
	doneInitialValue := false
	m.done = &doneInitialValue

	// Preserve the current message value
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Key("message").
				Title("Describe Your Policy").Value(m.message),
			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Value(m.done).
				Affirmative("Yes").
				Negative("Refresh"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
}
