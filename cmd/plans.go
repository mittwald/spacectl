package cmd

import (
	"github.com/spf13/cobra"
)

var plansCmd = &cobra.Command{
	Use:     "plans",
	Aliases: []string{"plans"},
	Short:   "View available plans",
	Long:    `This command allows you to view available plans`,
}

func init() {
	RootCmd.AddCommand(plansCmd)
}
