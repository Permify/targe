package main

import (
	"os"

	"github.com/Permify/kivo/pkg/cmd"
	"github.com/Permify/kivo/pkg/cmd/aws"
)

func main() {
	root := cmd.NewRootCommand()

	awsCommand := aws.NewAwsCommand()
	root.AddCommand(awsCommand)

	// aws_iam_actions.SaveActions()

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
