package cmd

import (
	"errors"
	"fmt"
	"github.com/mittwald/spacectl/client/backups"
	"github.com/spf13/cobra"
)

var backupsRecoverFlags struct {
	Force       bool
	Files       []string
	Stage       string
	NoFiles     bool
	NoDatabases bool
	NoMetadata  bool
}

var backupsRecoverCmd = &cobra.Command{
	Use:   "recover <backup-id> [--stage=<stage>] [--without-databases] [--without-metadata] [--without-files|--file=<file>...]",
	Short: "Recover a specific backup",
	Example: `  Recover one specific file:
    spacectl backup recover ae46198c-7d69-44c8-8670-1968703f4aaf --without-databases --without-metadata --file=/path/to/file.foo

  Full recovery:
    spacectl backup recover ae46198c-7d69-44c8-8670-1968703f4aaf`,
	Long: `This command recovers a specific backup.

You can use several command line switches to control what should be recovered.
By default, this command will trigger a recovery process that recovers
*everything* from the specified backup. This will include:

  - all files from your file system
  - the entire database content
  - meta data like the installed software version, configured HTTP hosts,
    cron jobs and more

To disable the recovery of certain items, you can use the --without-* flags
listed below.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.SilenceUsage = false
			return errors.New("missing argument: backup ID")
		}

		backup, err := api.Backups().Get(args[0])
		if err != nil {
			return err
		}

		stage := backupsRecoverFlags.Stage
		if stage == "" {
			stage = backup.Stage.Name
		}

		fileRecoverySpec := backups.RecoverySpec{Type: backups.RecoverAll}
		databaseRecoverySpec := backups.RecoverySpec{Type: backups.RecoverAll}
		metadataRecoverySpec := backups.RecoverySpec{Type: backups.RecoverAll}

		if backupsRecoverFlags.NoFiles {
			fileRecoverySpec.Type = backups.RecoverNone
		} else if len(backupsRecoverFlags.Files) > 0 {
			fileRecoverySpec.Type = backups.RecoverSpecific
			fileRecoverySpec.Items = backupsRecoverFlags.Files
		}

		if backupsRecoverFlags.NoDatabases {
			databaseRecoverySpec.Type = backups.RecoverNone
		}

		if backupsRecoverFlags.NoMetadata {
			metadataRecoverySpec.Type = backups.RecoverNone
		}

		recovery, err := api.Backups().Recover(backup.ID, stage, fileRecoverySpec, databaseRecoverySpec, metadataRecoverySpec)
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
	backupsRecoverCmd.Flags().StringVarP(&backupsRecoverFlags.Stage, "stage", "e", "", "Target stage for backup recovery")
	backupsRecoverCmd.Flags().BoolVar(&backupsRecoverFlags.NoFiles, "without-files", false, "Set to disable file recovery")
	backupsRecoverCmd.Flags().BoolVar(&backupsRecoverFlags.NoDatabases, "without-databases", false, "Set to disable database recovery")
	backupsRecoverCmd.Flags().BoolVar(&backupsRecoverFlags.NoMetadata, "without-metadata", false, "Set to disable metadata recovery")
}
