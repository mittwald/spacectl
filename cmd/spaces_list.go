package cmd

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
	"github.com/mittwald/spacectl/cmd/helper"
)

// listCmd represents the list command
var spacesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List Spaces",
	Long:    `Lists spaces`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var ownedSpaces []spaces.Space
		var err error
		teamID := viper.GetString("teamID")

		if teamID != "" {
			ownedSpaces, err = api.Spaces().ListByTeam(teamID)
		} else {
			ownedSpaces, err = api.Spaces().List()
		}

		if err != nil {
			return err
		}

		if len(ownedSpaces) == 0 {
			fmt.Println("No Spaces found.")
			return nil
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("ID", "DNS LABEL", "TEAM", "NAME", "STAGES", "RUNNING", "CREATED")

		for _, space := range ownedSpaces {
			since := helper.HumanReadableDateDiff(time.Now(), space.CreatedAt)

			running := "no"
			if space.Running {
				running = "yes"
			}

			table.AddRow(
				space.ID,
				space.Name.DNSName,
				space.Team.Name,
				space.Name.HumanReadableName,
				strings.Join(space.StagesNames(), ", "),
				running,
				since+" ago",
			)
		}

		fmt.Println(table)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
