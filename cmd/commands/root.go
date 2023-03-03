package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "commit-gpt",
		Short: "The commit-gpt command line tool",
		Long:  `The commit-gpt tool is used by developers to generate the commit messages using AI`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your command '%s'", err)
		os.Exit(1)
	}
}
