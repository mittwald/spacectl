package cmd

import (
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"os"
)

var appsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List applications",
	Long:    `This command allows you to view applications that are offered as a managed service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		softwares, err := api.Applications().ListWithVersions()
		if err != nil {
			return err
		}

		v := view.TabularAppsView{}
		v.List(softwares, os.Stdout)

		return nil
	},
}

func init() {
	appsCmd.AddCommand(appsListCmd)
}
