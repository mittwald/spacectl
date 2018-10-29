package cmd

import (
	"fmt"

	"errors"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var teamsRolesCmd = &cobra.Command{
	Use:   "roles -t <team-id>",
	Short: "List team roles",
	Long:  `List roles that members may have in a given team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID := viper.GetString("teamID")
		if teamID == "" {
			return errors.New("must provide team ID (--team-id, -t or $SPACES_TEAM_ID variable)")
		}

		roles, err := api.Teams().ListRoles(teamID)
		if err != nil {
			return err
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("ID")

		for _, m := range roles {
			table.AddRow(m.Identifier)
		}

		fmt.Println(table)
		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamsRolesCmd)
}
