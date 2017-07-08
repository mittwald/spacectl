package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"errors"
	"strings"
	"regexp"
	"github.com/spf13/viper"
)

var spaceInitForce bool
var spaceInitName string
var spaceInitLabel string

var spacesInitCmd = &cobra.Command{
	Use:   "init -t <team-name>",
	Short: "Initialize a new space",
	Long: `This command initializes a new space.

Note that this command does not actually do anything, except creating a new
Spacefile in your current working directory. You can then use the "spacectl space apply"
command to actually apply the declaration within the Spacefile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		teamID := viper.GetString("teamID")

		if teamID != "" {
			return errors.New("must provide team (--team, -t or $SPACES_TEAM_ID)")
		}

		if spaceInitName == "" {
			return errors.New("must provide name (--name or -n)")
		}

		if spaceInitLabel == "" {
			spaceInitLabel = regexp.MustCompile("[^a-z0-9-]").ReplaceAllString(strings.ToLower(spaceInitName), "-")
			fmt.Printf("Using '%s' as auto-generated DNS label\n", spaceInitLabel)
		}

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesInitCmd)

	spacesInitCmd.Flags().BoolVarP(&spaceInitForce, "force", "f", false, "Override existing Spacefile without asking")
	spacesInitCmd.Flags().StringVarP(&spaceInitName, "name", "n", "", "Name of the new Space")
	spacesInitCmd.Flags().StringVarP(&spaceInitLabel, "dns-label", "l", "", "DNS label of the new Space. Must be unique within the team.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spacesInitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spacesInitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
