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
	"github.com/mittwald/spacectl/client/teams"
	"github.com/gosuri/uitable"
)

// inviteCmd represents the invite command
var teamInviteCmd = &cobra.Command{
	Use:   "invite -t <team-id> -e <email> -m <message>",
	Short: "Invite new users to your team",
	Long: `Invite a new user into your team`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var invite teams.Invite

		teamID := cmd.Flag("team-id").Value.String()
		email := cmd.Flag("email").Value.String()
		userID := cmd.Flag("user-id").Value.String()
		message := cmd.Flag("message").Value.String()
		role := cmd.Flag("role").Value.String()

		if teamID == "" {
			return errors.New("must provide team (--team-id or -t)")
		}

		if email == "" && userID == "" {
			return errors.New("must provide user (either --email|-e or --user-id|-u)")
		}

		if message == "" {
			return errors.New("must provide message (--message or -m)")
		}

		if email != "" {
			fmt.Printf("inviting user \"%s\" into team %s\n", email, teamID)
			invite, err = api.Teams().InviteByEmail(teamID, email, message, role)
		}

		if err != nil {
			return err
		}

		fmt.Printf("invite %s issued\n", invite.ID)

		table := uitable.New()
		table.MaxColWidth = 80
		table.Wrap = true

		table.AddRow("ID:", invite.ID)
		table.AddRow("Message:", invite.Message)
		table.AddRow("State:", invite.State)

		fmt.Println(table)

		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamInviteCmd)

	teamInviteCmd.Flags().StringP("email", "e", "", "Email address of the user to invite");
	teamInviteCmd.Flags().StringP("user-id", "u", "", "User ID of the user to invite")
	teamInviteCmd.Flags().StringP("message", "m", "", "Invitation message")
	teamInviteCmd.Flags().StringP("role", "r", "", "User role")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inviteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inviteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
