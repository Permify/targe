package ai

import (
	"github.com/spf13/cobra"

	"github.com/Permify/kivo/internal/config"
)

// NewAICommand -
func NewAICommand(cfg *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "ai",
		Short: "",
		Long:  ``,
	}

	return command
}
