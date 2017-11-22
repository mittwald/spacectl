package cmd

import (
	"github.com/mittwald/spacectl/client/backups"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/view"
	"os"
)

var backupsCreateFlags struct {
	SpaceID     string
	StageName   string
	Keep        bool
	Description string
}

var backupsCreateCmd = &cobra.Command{
	Use:   "create --space <space-id> --stage <stage> [--keep|-k] [--description|-d <description>]",
	Short: "Create new backup",
	Long:  `Creates a new backup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &backupsCreateFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		f := &backupsCreateFlags
		backup, err := api.Backups().Create(
			space.ID,
			f.StageName,
			f.Keep,
			f.Description,
		)

		if err != nil {
			return err
		}

		v := view.TabularBackupView{}
		v.Detail(backup, []backups.Recovery{}, space, os.Stdout)

		return nil
	},
}

func init() {
	backupsCmd.AddCommand(backupsCreateCmd)

	f := backupsCreateCmd.Flags()

	f.StringVarP(&backupsCreateFlags.SpaceID, "space", "s", "", "Space ID")
	f.StringVarP(&backupsCreateFlags.StageName, "stage", "e", "", "Stage name")
	f.BoolVarP(&backupsCreateFlags.Keep, "keep", "k", false, "Create backup that is kept indefinitely")
	f.StringVarP(&backupsCreateFlags.Description, "description", "d", "", "Description for the new backup (optional)")
}
