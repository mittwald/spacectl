package cmd

import (
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:     "keys",
	Aliases: []string{"key"},
	Short:   "Manage SSH public keys",
	Long:    `This command allows you to manage SSH public keys`,
}

func init() {
	RootCmd.AddCommand(keysCmd)
}
