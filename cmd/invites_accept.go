package cmd

import (
	"github.com/spf13/cobra"
	"errors"
	"fmt"
)

var acceptInviteID string

// acceptCmd represents the accept command
var invitesAcceptCmd = &cobra.Command{
	Use:   "accept",
	Short: "Accepts a pending invitation",
	Long: `This command accepts a pending invitation and makes you a member of the team you were invited into.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if acceptInviteID == "" {
			return errors.New("must specify invite ID (--invite or -i)")
		}

		err := api.Invites().Accept(acceptInviteID)
		if err != nil {
			return err
		}

		fmt.Println("Invite accepted")

		return nil
	},
}

func init() {
	invitesCmd.AddCommand(invitesAcceptCmd)

	invitesAcceptCmd.Flags().StringVarP(&acceptInviteID, "invite", "i", "", "ID of the invite to accept")
}
