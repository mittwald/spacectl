package cmd

import (
	"github.com/spf13/cobra"
	"errors"
	"github.com/mittwald/spacectl/view"
	"os"
)

var keysShowCmd = &cobra.Command{
	Use:     "get <key-id>",
	Short:   "Get specific SSH key",
	Long:    `This command gets a specific SSH key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			RootCmd.SilenceUsage = false
			return errors.New("Missing argument: Key ID")
		}

		keyID := args[0]
		key, err := api.SSHKeys().Get(keyID)
		if err != nil {
			return err
		}

		v := view.TabularKeyDetailView{}
		v.KeyDetail(key, os.Stdout)

		return nil
	},
}

func init() {
	keysCmd.AddCommand(keysShowCmd)
}
