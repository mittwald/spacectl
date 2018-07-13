package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var revokeInviteID string

// acceptCmd represents the accept command
var invitesRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revokes a pending invitation",
	Long:  `This command revokes a pending invitation, stopping the invited user from becoming a member of the team`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if revokeInviteID == "" {
			return errors.New("must specify invite ID (--invite or -i)")
		}

		err := api.Invites().Revoke(revokeInviteID)
		if err != nil {
			return err
		}

		fmt.Println("Invite revoked")

		return nil
	},
}

func init() {
	invitesCmd.AddCommand(invitesRevokeCmd)

	invitesRevokeCmd.Flags().StringVarP(&revokeInviteID, "invite", "i", "", "ID of the invite to revoke")
}
