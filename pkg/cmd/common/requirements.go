package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Permify/kivo/internal/requirements"
)

type RequirementsManager struct {
	requirements []requirements.Requirement
	index        int
	width        int
	height       int
	spinner      spinner.Model
	progress     progress.Model
	done         bool
}

var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func NewRequirements() RequirementsManager {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return RequirementsManager{
		requirements: requirements.GetRequirements(),
		spinner:      s,
		progress:     p,
	}
}

func (m RequirementsManager) Init() tea.Cmd {
	return tea.Batch(downloadAndInstall(m.requirements[m.index]), m.spinner.Tick)
}

func (m RequirementsManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case installedPkgMsg:
		pkg := m.requirements[m.index]
		if m.index >= len(m.requirements)-1 {
			// Everything's been installed. We're done!
			m.done = true
			return m, tea.Sequence(
				tea.Printf("%s %s", checkMark, pkg.GetName()), // print the last success message
				tea.Quit, // exit the program
			)
		}

		// Update progress bar
		m.index++
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.requirements)))

		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, pkg),         // print success message above our program
			downloadAndInstall(m.requirements[m.index]), // download the next package
		)
	case installErrorMsg:
		// Update state for errors
		tea.Printf("Error installing package: %v", msg.Err)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m RequirementsManager) View() string {
	n := len(m.requirements)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! Installed %d requirements.\n", n))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	pkgName := currentPkgNameStyle.Render(m.requirements[m.index].GetName())
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Installing " + pkgName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + pkgCount
}

// Message types
type installedPkgMsg struct {
	Name string
}

type installErrorMsg struct {
	Err error
}

// downloadAndInstall asynchronously downloads and installs a requirement
func downloadAndInstall(requirement requirements.Requirement) tea.Cmd {
	return func() tea.Msg {
		// Simulate installation
		err := requirement.Install()
		if err != nil {
			return installErrorMsg{Err: err}
		}

		// Return a success message
		return installedPkgMsg{Name: requirement.GetName()}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
