package cmd

import (
	"github.com/spf13/cobra"
)

var softwareCmd = &cobra.Command{
	Use:     "software",
	Aliases: []string{"sw"},
	Short:   "View available managed software",
	Long:    `This command allows you to view software that is offered as a managed service`,
}

func init() {
	RootCmd.AddCommand(softwareCmd)
}
