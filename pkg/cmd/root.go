package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand - Creates new root command
func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kivo",
		Short: "",
		Long:  ``,
	}
}
