package main

import (
	"fmt"
	"os"
	"soccer-cli/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	command := rootCmd.CobraCommand()

	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
