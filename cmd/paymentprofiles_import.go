package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

var paymentProfilesImportCmd = &cobra.Command{
	Use:     "connect",
	Short:   "Connect existing profile",
	Long:    `This command connects an existing Mittwald customer number with a SPACES payment profile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Mittwald customer number: ")
		customerNumber, err := reader.ReadString('\n')

		if err != nil {
			return err
		}

		fmt.Print("Enter password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))

		profile, err := api.Payment().ConnectProfile(customerNumber, string(bytePassword))

		if err != nil {
			return err
		}

		fmt.Printf("successfully imported payment profile %s\n", profile.ID)

		return nil
	},
}

func init() {
	paymentProfilesCmd.AddCommand(paymentProfilesImportCmd)
}
