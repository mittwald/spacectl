package cmd

import (
	"errors"
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"os"
)

var backupsShowCmd = &cobra.Command{
	Use:   "show <backup-id>",
	Short: "Show details regarding a specific backup",
	Long:  `This command shows details regarding a specific backup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.SilenceUsage = false
			return errors.New("missing argument: backup ID")
		}

		backup, err := api.Backups().Get(args[0])
		if err != nil {
			return err
		}

		var space *spaces.Space
		if backup.Space != nil {
			space, err = api.Spaces().GetByID(backup.Space.ID)
			if err != nil {
				return err
			}
		}

		recoveries, err := api.Recoveries().ListForBackup(backup)
		if err != nil {
			return err
		}

		v := view.TabularBackupView{}
		v.Detail(backup, recoveries, space, os.Stdout)

		return nil
	},
}

func init() {
	backupsCmd.AddCommand(backupsShowCmd)
}
