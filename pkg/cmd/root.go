package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/ai"
	configc "github.com/Permify/kivo/pkg/cmd/config"

	"github.com/Permify/kivo/internal/config"
	"github.com/Permify/kivo/pkg/cmd/aws"
)

type RootModel struct {
	command  string
	choice   string
	quitting bool
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "", tea.KeyEnter.String():
			return RootModel{command: m.command, choice: "yes", quitting: true}, tea.Quit
		case "n", "N":
			return RootModel{command: m.command, choice: "no", quitting: true}, tea.Quit
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m RootModel) View() string {
	if m.quitting {
		if m.choice == "yes" {
			return ""
		}
		return fmt.Sprintf(
			"%s\n\n",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).Render("✘ Command aborted."),
		)
	}

	// Define styles
	brandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Padding(0, 0)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("6")).
		Underline(true)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Italic(true).
		PaddingLeft(2)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("3")).
		PaddingTop(1)

	// Render sections
	brand := brandStyle.Render("Generating your command...")
	header := headerStyle.Render("Here’s your command:")

	// Format the command
	formattedCommand := formatCommand(m.command, 2)
	message := messageStyle.Render(fmt.Sprintf("➤ kivo %s", formattedCommand))
	prompt := promptStyle.Render("Would you like to use this command? (Y/n):")

	// Combine output
	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s", brand, header, message, prompt)
}

// NewRootCommand - Creates new root command
func NewRootCommand() *cobra.Command {
	// Load the configuration
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Failed to load configuration:", err)
		os.Exit(1)
	}

	root := &cobra.Command{
		Use:   "kivo",
		Short: "",
		Long:  ``,
		RunE:  r(cfg),
	}

	f := root.Flags()

	f.String("m", "", "message")

	// SilenceUsage is set to true to suppress usage when an error occurs
	root.SilenceUsage = true

	root.PreRun = func(cmd *cobra.Command, args []string) {
		RegisterRootFlags(f)
	}

	configCommand := configc.NewConfigCommand()
	awsCommand := aws.NewAwsCommand(cfg)

	root.AddCommand(awsCommand, configCommand)

	return root
}

func r(cfg *config.Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		message := viper.GetString("m")

		gptResponse, err := ai.UserPrompt(cfg.OpenaiApiKey, message)
		if err != nil {
			return err
		}

		command := ai.GenerateCLICommand(gptResponse)

		// Bubble Tea program setup
		program := tea.NewProgram(&RootModel{command: command})
		mod, err := program.Run()
		if err != nil {
			return fmt.Errorf("program encountered an error: %w", err)
		}

		// Check user choice
		if result, ok := mod.(RootModel); ok && result.choice == "yes" {
			var args []string
			args = append(args, strings.Split(command, " ")...)
			cmd.SetArgs(args)
			return cmd.Root().Execute()
		}

		return nil
	}
}

// Helper function to format a command
func formatCommand(command string, maxWordsPerLine int) string {
	parts := strings.Fields(command) // Split command into words
	var result []string
	var line []string

	for _, part := range parts {
		line = append(line, part)
		if len(line) >= maxWordsPerLine {
			result = append(result, strings.Join(line, " ")+" \\")
			line = []string{}
		}
	}

	// Add the remaining words
	if len(line) > 0 {
		result = append(result, strings.Join(line, " "))
	}

	return strings.Join(result, "\n      ")
}
