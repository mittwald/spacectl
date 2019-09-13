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

	Stages                int64
	Storage               string
	Pods                  int64
	SkipTestingPeriod     bool
	BackupIntervalMinutes uint64
}

var spacesConnectCmd = &cobra.Command{
	Use:   "connect --plan <plan> --profile <profile-id> -s <space-id> [--stages <stages>] [--storage <storage>] [--pods <pods>] [--skip-testing-period] [--backup-interval-minutes <interval>]",
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

		if spaceConnectFlags.Stages >= 0 {
			opts = append(opts, spaces.WithStages(uint64(spaceConnectFlags.Stages)))
		}

		if spaceConnectFlags.Pods >= 0 {
			opts = append(opts, spaces.WithPods(uint64(spaceConnectFlags.Pods)))
		}

		if spaceConnectFlags.SkipTestingPeriod {
			opts = append(opts, spaces.WithoutTestingPeriod())
		}

		if spaceConnectFlags.BackupIntervalMinutes != 0 {
			opts = append(opts, spaces.WithBackupIntervalMinutes(spaceConnectFlags.BackupIntervalMinutes))
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

	flags := spacesConnectCmd.Flags()

	flags.StringVarP(&spaceConnectFlags.SpaceID, "space", "s", "", "ID of the space to connect")
	flags.StringVar(&spaceConnectFlags.PlanID, "plan", "", "Plan ID to use")
	flags.StringVar(&spaceConnectFlags.PaymentProfileID, "profile", "", "Payment profile ID")

	flags.Int64Var(&spaceConnectFlags.Stages, "stages", -1, "amount of preprovisioned stages (CAUTION: Additional charges will apply!)")
	flags.StringVar(&spaceConnectFlags.Storage, "storage", "10G", "amount of preprovisioned storage (CAUTION: Additional charges will apply!)")
	flags.Int64Var(&spaceConnectFlags.Pods, "pods", -1, "amount of preprovisioned pods (CAUTION: Additional charges will apply!)")
	flags.BoolVar(&spaceConnectFlags.SkipTestingPeriod, "skip-testing-period", false, "skip testing period")
	flags.Uint64Var(&spaceConnectFlags.BackupIntervalMinutes, "backup-interval-minutes", 0, "desired minimum backup interval in minutes (CAUTION: Additional charges will apply!)")

	_ = spacesConnectCmd.MarkFlagRequired("plan")
	_ = spacesConnectCmd.MarkFlagRequired("profile")
}
