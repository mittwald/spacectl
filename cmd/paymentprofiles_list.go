package cmd

import (
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"os"
)

var paymentProfilesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List payment profiles",
	Long:    `Lists payment profiles`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := api.Payment().ListPaymentProfiles()

		if err != nil {
			return err
		}

		v := view.TabularPaymentProfileListView{}
		v.List(profiles, os.Stdout)

		return nil
	},
}

func init() {
	paymentProfilesCmd.AddCommand(paymentProfilesListCmd)
}
