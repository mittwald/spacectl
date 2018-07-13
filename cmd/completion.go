package cmd

import (
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates a Bash completion script",
	Long:  "Generates a Bash completion script",
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
