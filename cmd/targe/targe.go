package main

import (
	"os"

	"github.com/Permify/targe/pkg/cmd"
)

func main() {
	root := cmd.NewRootCommand()

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
