package groups

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Permify/kivo/pkg/cmd/ai"
)

type CreatePolicy struct {
	controller  *Controller
	message     string
	senderStyle lipgloss.Style
	viewport    viewport.Model
	textarea    textarea.Model
	err         error
}

func NewCreatePolicy(controller *Controller) CreatePolicy {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 15)
	vp.SetContent(`Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return CreatePolicy{
		controller:  controller,
		viewport:    vp,
		textarea:    ta,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func (m CreatePolicy) Init() tea.Cmd {
	return textarea.Blink
}

func (m CreatePolicy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":

			if m.controller.State.policy != nil {
				return Switch(m.controller.Next(), 0, 0)
			} else {

				m.message = m.textarea.Value()

				policy, err := ai.GeneratePolicy("", m.message)
				if err != nil {
					m.err = err
				}

				policyJson, err := json.Marshal(policy)
				if err != nil {
					m.err = err
				}

				m.viewport.SetContent(string(policyJson))
				m.textarea.Reset()
				m.viewport.GotoBottom()

				//var result map[string]interface{}
				//m.err = json.Unmarshal([]byte(jsonStr), &result)
				//m.state.SetPolicy(&Policy{
				//	Arn:      "new",
				//	Name:     "New",
				//	Document: result,
				//})
			}
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m CreatePolicy) View() string {
	if m.err != nil {
		return createPolicyStyle.Render(m.err.Error())
	}

	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
