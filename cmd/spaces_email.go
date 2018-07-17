package cmd

import (
	"github.com/spf13/cobra"
)

var spaceEmailCmd = &cobra.Command{
	Use:     "emails",
	Aliases: []string{"email"},
	Short:   "Manage outgoing emails",
	Long:    `This command allows you to manage outgoing emails`,
}

func init() {
	spacesCmd.AddCommand(spaceEmailCmd)
}
