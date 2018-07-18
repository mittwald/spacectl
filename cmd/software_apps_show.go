package cmd

import (
	"errors"
	"os"

	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var appsShowFlags struct {
	AppID string
}

var appsShowCmd = &cobra.Command{
	Use:   "show <app-id>",
	Short: "show application",
	Long:  `This command allows you to inspect an application's versions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		appArg := appsShowFlags.AppID
		if appArg == "" && len(args) >= 1 {
			appArg = args[0]
		}

		if appArg == "" {
			RootCmd.SilenceUsage = false
			return errors.New("Missing argument: Application ID")
		}

		app, err := api.Applications().Get(appArg)
		if err != nil {
			return err
		}

		v := view.TabularAppView{}
		v.List(app.Identifier, app.Name, app.Versions, os.Stdout)

		return nil
	},
}

func init() {
	appsCmd.AddCommand(appsShowCmd)
	appsShowCmd.Flags().StringVarP(&appsShowFlags.AppID, "app", "a", "", "Application Identifier")
}
