package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/cmd/helper"
	"time"
)

var teamsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List teams",
	Long:    `Lists all teams that you have access to.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		teams, err := api.Teams().List()
		if err != nil {
			return err
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("ID", "DNS LABEL", "NAME", "CREATED")

		for _, team := range teams {
			table.AddRow(team.ID, team.DNSName, team.Name, helper.HumanReadableDateDiff(time.Now(), team.CreatedAt) + " ago")
		}

		fmt.Println(table)

		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamsListCmd)
}
