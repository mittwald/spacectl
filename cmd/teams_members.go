// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"errors"
	"github.com/gosuri/uitable"
	"strings"
)

// membersCmd represents the members command
var teamsMembersCmd = &cobra.Command{
	Use:   "members -t <team-id>",
	Short: "List team members",
	Long: `List users that are members of a given team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID := cmd.Flag("team-id").Value.String()

		if teamID == "" {
			return errors.New("must provide team ID (--team-id or -t)")
		}

		memberships, err := spaces.Teams().ListMembers(teamID)
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

	teamsMembersCmd.Flags().StringP("team-id", "t", "", "Team ID to list members for")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// membersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// membersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
