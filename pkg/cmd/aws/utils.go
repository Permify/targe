package aws

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func RootModel(m tea.Model) Users {
	return Users{
		model: m,
	}
}

func folderExists(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Helper function to extract the resource name from the ARN
func parseResourceNameFromArn(arn string) string {
	parts := strings.Split(arn, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1] // Return the last part of the ARN
	}
	return arn
}
