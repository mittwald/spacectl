package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var teamsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List teams",
	Long:  `Lists all teams that you have access to.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		teams, err := spaces.Teams().List()
		if err != nil {
			return err
		}

		fmt.Println("list called")
		fmt.Println(teams)

		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamsListCmd)
}
