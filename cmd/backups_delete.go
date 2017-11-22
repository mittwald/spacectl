package cmd

import (
	"errors"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"bytes"
	"fmt"
)

var backupsDeleteFlags struct {
	Force bool
}

var backupsDeleteCmd = &cobra.Command{
	Use:     "delete [--yes|-y] <backup-id>",
	Aliases: []string{"rm"},
	Short:   "Delete a specific backup",
	Long:    `This command deletes a specific backup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.SilenceUsage = false
			return errors.New("missing argument: backup ID")
		}

		backup, err := api.Backups().Get(args[0])
		if err != nil {
			return err
		}

		recoveries, err := api.Recoveries().ListForBackup(backup)
		if err != nil {
			return err
		}

		space, err := api.Spaces().GetByID(backup.Space.ID)
		if err != nil {
			return err
		}

		if !backupsDeleteFlags.Force {
			buf := bytes.Buffer{}
			view.TabularBackupView{}.Detail(backup, recoveries, space, &buf)

			ok, _ := view.Confirm("Once this backup is deleted, you will NOT be able to get it back.", buf.String())
			if !ok {
				fmt.Println("Aborting backup deletion.")
				return nil
			}
		}

		err = api.Backups().Delete(backup.ID)
		if err != nil {
			return err
		}

		fmt.Println("Backup deleted.")

		return nil
	},
}

func init() {
	backupsCmd.AddCommand(backupsDeleteCmd)

	backupsDeleteCmd.Flags().BoolVarP(&backupsDeleteFlags.Force, "yes", "y", false, "Do not prompt for confirmation")
}
