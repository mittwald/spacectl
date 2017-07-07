package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"errors"
	"github.com/gosuri/uitable"
	"strings"
	"github.com/spf13/viper"
)

// membersCmd represents the members command
var teamsMembersCmd = &cobra.Command{
	Use:   "members -t <team-id>",
	Short: "List team members",
	Long: `List users that are members of a given team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID := viper.GetString("teamID")
		if teamID == "" {
			return errors.New("must provide team ID (--team-id, -t or $SPACES_TEAM_ID variable)")
		}

		memberships, err := api.Teams().ListMembers(teamID)
		if err != nil {
			return err
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("ID", "ROLE", "NAME", "E-MAIL")

		for _, m := range memberships {
			name := strings.TrimSpace(fmt.Sprintf("%s %s", m.User.FirstName, m.User.LastName))
			table.AddRow(m.User.ID, m.Role, name, m.User.Email)
		}

		fmt.Println(table)
		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamsMembersCmd)
}
