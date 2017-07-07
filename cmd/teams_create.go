// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/gosuri/uitable"
	"errors"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create -n <team-name> -d <dns-label>",
	Short: "Create a new team",
	Long: `Creates a new team. Afterwards, you will have "Owner" access on the newly created team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := cmd.Flag("name").Value.String()
		dnsLabel := cmd.Flag("dns-label").Value.String()

		if name == "" {
			return errors.New("must provide name (--name or -n)")
		}
		if dnsLabel == "" {
			return errors.New("must provide DNS label (--dns-label or -d)")
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
	createCmd.Flags().StringP("dns-label", "d", "", "The new team's DNS label")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
