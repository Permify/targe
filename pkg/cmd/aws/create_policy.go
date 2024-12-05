package aws

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// var createPolicyStyle = lipgloss.NewStyle().Margin(1, 2)
type CreatePolicyModel struct {
	user        User
	response    string
	resource    Resource
	senderStyle lipgloss.Style
	viewport    viewport.Model
	textarea    textarea.Model
}

func CreatePolicy(user User, resource Resource) CreatePolicyModel {
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

	return CreatePolicyModel{
		user:        user,
		resource:    resource,
		viewport:    vp,
		textarea:    ta,
		response:    "",
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func (m CreatePolicyModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m CreatePolicyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:

			// m.textarea.Value()

			m.viewport.SetContent(`
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": "*"
        }
    ]
}
			`)
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m CreatePolicyModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
