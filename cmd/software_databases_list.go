package cmd

import (
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"os"
)

var dbsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List database engines",
	Long:    `This command allows you to view database engines that are offered as a managed service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		softwares, err := api.Databases().ListWithVersions()
		if err != nil {
			return err
		}

		v := view.TabularAppsView{}
		v.List(softwares, os.Stdout)

		return nil
	},
}

func init() {
	dbsCmd.AddCommand(dbsListCmd)
}
