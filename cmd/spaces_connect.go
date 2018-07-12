package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/cmd/helper"
)

var spaceConnectFlags struct {
	SpaceID string
	PlanID string
	PaymentProfileID string
}

var spacesConnectCmd = &cobra.Command{
	Use:   "connect --plan <plan> --profile <profile-id> -s <space-id>",
	Short: "Connect a Space with a payment profile and a plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spaceConnectFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		_, err = api.Spaces().ConnectWithPaymentProfile(space.ID, spaceConnectFlags.PaymentProfileID, spaceConnectFlags.PlanID)
		if err != nil {
			return err
		}

		fmt.Printf("Space connected to payment profile '%s' with plan '%s'\n", spaceConnectFlags.PaymentProfileID, spaceConnectFlags.PlanID)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesConnectCmd)

	spacesConnectCmd.Flags().StringVarP(&spaceConnectFlags.SpaceID, "space", "s", "", "ID of the space to connect")
	spacesConnectCmd.Flags().StringVar(&spaceConnectFlags.PlanID, "plan", "", "Plan ID to use")
	spacesConnectCmd.Flags().StringVar(&spaceConnectFlags.PaymentProfileID, "profile", "", "Payment profile ID")

	spacesConnectCmd.MarkFlagRequired("plan")
	spacesConnectCmd.MarkFlagRequired("profile")
}
