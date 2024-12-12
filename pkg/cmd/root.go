package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/config"
	"github.com/Permify/kivo/pkg/cmd/ai"
	"github.com/Permify/kivo/pkg/cmd/aws"
)

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
	}

	awsCommand := aws.NewAwsCommand(cfg)
	aiCommand := ai.NewAICommand(cfg)

	root.AddCommand(awsCommand)
	root.AddCommand(aiCommand)

	return root
}
