package cmd

import (
	"fmt"

	"crypto/md5"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/spf13/cobra"
	"time"
)

var keysListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List SSH keys",
	Long:    `This command lists SSH keys that are already managed for you.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := api.SSHKeys().List()
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			fmt.Println("No SSH keys found.")
			return nil
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("ID", "TYPE", "FINGERPRINT", "CREATED")

		for _, key := range keys {
			fp := md5.Sum(key.Key)

			table.AddRow(
				key.ID,
				key.CipherAlgorithm,
				fmt.Sprintf("%x", fp),
				helper.HumanReadableDateDiff(time.Now(), key.CreatedAt)+" ago",
			)
		}

		fmt.Println(table)

		return nil
	},
}

func init() {
	keysCmd.AddCommand(keysListCmd)
}
