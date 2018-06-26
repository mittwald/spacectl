package cmd

import (
	"github.com/spf13/cobra"
)

var dbsCmd = &cobra.Command{
	Use:     "databases",
	Aliases: []string{"dbs", "db"},
	Short:   "View available managed database engines",
	Long:    `This command allows you to view database engines that are offered as a managed service`,
}

func init() {
	softwareCmd.AddCommand(dbsCmd)
}
