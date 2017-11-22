package cmd

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"time"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/client/backups"
)

var backupsListFlags struct {
	SpaceID string
	StageName string
	Keep bool
}

var backupsListCmd = &cobra.Command{
	Use:     "list",
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

		var b []backups.Backup
		if backupsListFlags.StageName != "" {
			b, err = api.Backups().ListForStage(space.ID, backupsListFlags.StageName, &opts)
		} else {
			b, err = api.Backups().ListForSpace(space.ID, &opts)
		}

		if err != nil {
			return err
		}

		if len(b) == 0 {
			fmt.Println("No backups found.")
			return nil
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("ID", "STAGE", "STATUS", "KEEP", "DESCRIPTION", "CREATED")

		for _, backup := range b {
			since := helper.HumanReadableDateDiff(time.Now(), backup.StartedAt)
			stage := backupsListFlags.StageName

			keep := "no"
			if backup.Keep {
				keep = "yes"
			}

			if stage == "" && backup.Stage != nil {
				stage = backup.Stage.Name
			}

			table.AddRow(
				backup.ID,
				stage,
				backup.Status,
				keep,
				backup.Description,
				since+" ago",
			)
		}

		fmt.Println(table)

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
