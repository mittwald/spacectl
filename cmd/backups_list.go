package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/client/backups"
	"github.com/mittwald/spacectl/view"
	"os"
)

var backupsListFlags struct {
	SpaceID string
	StageName string
	Keep bool
}

var backupsListCmd = &cobra.Command{
	Use:     "list --space <space-id> [--stage <stage>] [--keep|-k]",
	Aliases: []string{"ls"},
	Short:   "List backups",
	Long:    `Lists backups`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &backupsListFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		opts := backups.BackupListOptions{
			OnlyKeep: backupsListFlags.Keep,
		}

		var backupList []backups.Backup
		if backupsListFlags.StageName != "" {
			backupList, err = api.Backups().ListForStage(space.ID, backupsListFlags.StageName, &opts)
		} else {
			backupList, err = api.Backups().ListForSpace(space.ID, &opts)
		}

		if err != nil {
			return err
		}

		v := view.TabularBackupView{}
		v.List(backupList, backupsListFlags.StageName, os.Stdout)

		return nil
	},
}

func init() {
	backupsCmd.AddCommand(backupsListCmd)

	f := backupsListCmd.Flags()

	f.StringVarP(&backupsListFlags.SpaceID, "space", "s", "", "Space ID")
	f.StringVarP(&backupsListFlags.StageName, "stage", "e", "", "Stage name")
	f.BoolVarP(&backupsListFlags.Keep, "keep", "k", false, "Show only backups that are kept indefinitely")
}
