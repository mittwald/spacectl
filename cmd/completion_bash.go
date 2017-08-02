package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates auto-complete definitions for Bash",
	Long: "Generates auto-complete definitions for Bash",
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(completionBashCmd)
}
