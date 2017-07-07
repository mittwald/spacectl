package cmd

import (
	"github.com/spf13/cobra"
)

var spacesCmd = &cobra.Command{
	Use:     "spaces",
	Aliases: []string{"space", "spc", "s"},
	Short:   "Manage Spaces",
	Long:    `This command can be used to manage Spaces`,
}

func init() {
	RootCmd.AddCommand(spacesCmd)
}
