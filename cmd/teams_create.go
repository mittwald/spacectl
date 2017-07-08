package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/gosuri/uitable"
	"errors"
	"github.com/hashicorp/go-multierror"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create -n <team-name> -l <dns-label>",
	Short: "Create a new team",
	Long: `Creates a new team. Afterwards, you will have "Owner" access on the newly created team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var merr *multierror.Error

		name := cmd.Flag("name").Value.String()
		dnsLabel := cmd.Flag("dns-label").Value.String()

		if name == "" {
			merr = multierror.Append(merr, errors.New("Must provide name (--name or -n)"))
		}
		if dnsLabel == "" {
			merr = multierror.Append(merr, errors.New("Must provide DNS label(--dns-label or -l)"))
		}

		if merr.ErrorOrNil() != nil {
			RootCmd.SilenceUsage = false
			return merr
		}

		fmt.Printf("creating team '%s'\n", name)
		team, err := api.Teams().Create(name, dnsLabel)
		if err != nil {
			return err
		}

		fmt.Print("team successfully created\n\n")

		table := uitable.New()
		table.MaxColWidth = 80
		table.Wrap = true

		table.AddRow("ID:", team.ID)
		table.AddRow("Name:", team.Name)
		table.AddRow("DNS Label:", team.DNSName)

		fmt.Println(table)

		return nil
	},
}

func init() {
	teamsCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("name", "n", "", "The new team's name")
	createCmd.Flags().StringP("dns-label", "l", "", "The new team's DNS label")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
