package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/mittwald/spacectl/client/backups"
)

var backupsRecoverFlags struct {
	Force bool
	Files []string
	NoFiles bool
	NoDatabases bool
}

var backupsRecoverCmd = &cobra.Command{
	Use:     "recover <backup-id> [--without-databases] [--without-files|--file=<file>...]",
	Short:   "Recover a specific backup",
	Long:    `This command recovers a specific backup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.SilenceUsage = false
			return errors.New("missing argument: backup ID")
		}

		backup, err := api.Backups().Get(args[0])
		if err != nil {
			return err
		}

		fileRecoverySpec := backups.RecoverySpec{Type: backups.RecoverAll}
		databaseRecoverySpec := backups.RecoverySpec{Type: backups.RecoverAll}

		if backupsRecoverFlags.NoFiles {
			fileRecoverySpec.Type = backups.RecoverNone
		} else if len(backupsRecoverFlags.Files) > 0 {
			fileRecoverySpec.Type = backups.RecoverSpecific
			fileRecoverySpec.Items = backupsRecoverFlags.Files
		}

		if backupsRecoverFlags.NoDatabases {
			databaseRecoverySpec.Type = backups.RecoverNone
		}

		recovery, err := api.Backups().Recover(backup.ID, fileRecoverySpec, databaseRecoverySpec)
		if err != nil {
			return err
		}

		fmt.Printf("Recovery process started: %s\n", recovery.ID)

		return nil
	},
}

func init() {
	backupsCmd.AddCommand(backupsRecoverCmd)

	backupsRecoverCmd.Flags().StringSliceVarP(&backupsRecoverFlags.Files, "file", "r", []string{}, "List of files to recover")
	backupsRecoverCmd.Flags().BoolVar(&backupsRecoverFlags.NoFiles, "without-files", false, "Set to disable file recovery")
	backupsRecoverCmd.Flags().BoolVar(&backupsRecoverFlags.NoDatabases, "without-databases", false, "Set to disable database recovery")
}
