package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var keysDeleteCmd = &cobra.Command{
	Use:     "delete <key-id>",
	Aliases: []string{"rm"},
	Short:   "Delete specific SSH key",
	Long:    `This command deletes a specific SSH key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			RootCmd.SilenceUsage = false
			return errors.New("Missing argument: Key ID")
		}

		keyID := args[0]
		err := api.SSHKeys().Delete(keyID)
		if err != nil {
			return err
		}

		fmt.Printf("SSH key %s deleted.\n", keyID)

		return nil
	},
}

func init() {
	keysCmd.AddCommand(keysDeleteCmd)
}
