package cmd

import (
	"github.com/mittwald/spacectl/spacefile"
	"github.com/spf13/cobra"
)

var backupsCmd = &cobra.Command{
	Use:     "backups",
	Aliases: []string{"backup"},
	Short:   "Manage backups",
	Long:    `This command can be used to create and restore Space backups`,
}

func init() {
	RootCmd.AddCommand(backupsCmd)

	backupsCmd.PersistentFlags().StringVarP(&spaceFile, "spacefile", "f", "./"+spacefile.DefaultFilename, "Use Space defined in this file")
}
