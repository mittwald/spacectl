package cmd

import (
	"fmt"
	"github.com/mittwald/spacectl/client/spaces"

	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/spf13/cobra"
)

var spaceConnectFlags struct {
	SpaceID          string
	PlanID           string
	PaymentProfileID string

	Stages  uint64
	Storage string
	Pods    uint64
}

var spacesConnectCmd = &cobra.Command{
	Use:   "connect --plan <plan> --profile <profile-id> -s <space-id> --stages <stages> --storage <storage> --pods <pods>",
	Short: "Connect a Space with a payment profile and a plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spaceConnectFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		var opts []spaces.ConnectOption

		if spaceConnectFlags.Storage != "" {
			opts = append(opts, spaces.WithStorageStr(spaceConnectFlags.Storage))
		}

		if spaceConnectFlags.Stages != 0 {
			opts = append(opts, spaces.WithStages(spaceConnectFlags.Stages))
		}

		if spaceConnectFlags.Pods != 0 {
			opts = append(opts, spaces.WithPods(spaceConnectFlags.Pods))
		}

		_, err = api.Spaces().ConnectWithPaymentProfile(space.ID, spaceConnectFlags.PaymentProfileID, spaceConnectFlags.PlanID, opts...)
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

	spacesConnectCmd.Flags().Uint64Var(&spaceConnectFlags.Stages, "stages", 0, "amount of preprovisioned stages")
	spacesConnectCmd.Flags().StringVar(&spaceConnectFlags.Storage, "storage", "20G", "amount of preprovisioned storage")
	spacesConnectCmd.Flags().Uint64Var(&spaceConnectFlags.Pods, "pods", 0, "amount of preprovisioned pods")

	spacesConnectCmd.MarkFlagRequired("plan")
	spacesConnectCmd.MarkFlagRequired("profile")
}
