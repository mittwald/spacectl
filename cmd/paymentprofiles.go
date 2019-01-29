package cmd

import (
	"github.com/spf13/cobra"
)

var paymentProfilesCmd = &cobra.Command{
	Use:     "paymentprofiles",
	Aliases: []string{"paymentprofile", "pp"},
	Short:   "View accessible payment profiles",
	Long:    `This command lists payment profiles that you have access to`,
}

func init() {
	RootCmd.AddCommand(paymentProfilesCmd)
}
