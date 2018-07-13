package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/invites"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/spf13/cobra"
	"time"
)

var inviteListOutgoing bool = false

// listCmd represents the list command
var invitesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pending invites",
	Long:    `List pending invites`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var inv []invites.Invite
		var err error

		table := uitable.New()

		if inviteListOutgoing {
			inv, err = api.Invites().ListOutgoing()

			if err != nil {
				return err
			}

			if len(inv) == 0 {
				fmt.Println("No outgoing invites")
				return nil
			}

			table.AddRow("ID", "ISSUED", "STATE", "TEAM NAME", "INVITEE EMAIL", "MSG")

			for _, i := range inv {
				email := i.InviteeEmail
				if i.Invitee != nil {
					email = i.Invitee.Email
				}

				table.AddRow(
					i.ID,
					helper.HumanReadableDateDiff(time.Now(), i.IssuedAt)+" ago",
					i.State,
					i.Team.Name,
					email,
					i.Message,
				)
			}
		} else {
			inv, err = api.Invites().ListIncoming()

			if err != nil {
				return err
			}

			if len(inv) == 0 {
				fmt.Println("No incoming invites")
				return nil
			}

			table.AddRow("ID", "ISSUED", "STATE", "TEAM NAME", "INVITER EMAIL", "MSG")

			for _, i := range inv {
				table.AddRow(
					i.ID,
					helper.HumanReadableDateDiff(time.Now(), i.IssuedAt)+" ago",
					i.State,
					i.Team.Name,
					i.Inviter.Email,
					i.Message,
				)
			}
		}

		fmt.Println(table)

		return nil
	},
}

func init() {
	invitesCmd.AddCommand(invitesListCmd)

	invitesListCmd.Flags().BoolVar(&inviteListOutgoing, "out", false, "Set to list outgoing invites")
}
