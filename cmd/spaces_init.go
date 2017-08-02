package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"errors"
	"strings"
	"regexp"
	"github.com/spf13/viper"
	"github.com/hashicorp/go-multierror"
	"os"
	"github.com/mittwald/spacectl/spacefile"
	"github.com/fatih/color"
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
		var mErr *multierror.Error
		teamID := viper.GetString("teamID")

		if teamID == "" {
			mErr = multierror.Append(mErr, errors.New("must provide team (--team, -t or $SPACES_TEAM_ID)"))
		}

		if spaceInitName == "" {
			mErr = multierror.Append(mErr, errors.New("must provide name (--name or -n)"))
		}

		if mErr != nil {
			return mErr
		}

		if spaceInitLabel == "" {
			spaceInitLabel = regexp.MustCompile("[^a-z0-9-]").ReplaceAllString(strings.ToLower(spaceInitName), "-")
			fmt.Printf("Using %s as auto-generated DNS label\n", color.YellowString(spaceInitLabel))
		}

		filePath := "./" + spacefile.DefaultFilename

		_, err := os.Stat(filePath)
		if !os.IsNotExist(err) {
			if !spaceInitForce {
				RootCmd.SilenceUsage = false
				return fmt.Errorf(`The file '%s' already exists in the current directory.
Use the --force flag (or -f) to overwrite it.`, filePath)
			} else {
				fmt.Printf("Overwriting existing Spacefile at %s\n", color.YellowString(filePath))
			}
		}

		fh, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("Could not open '%s' for writing:\n    %s", filePath, err)
		}

		err = spacefile.Generate(teamID, spaceInitName, spaceInitLabel, fh)
		if err != nil {
			return fmt.Errorf("Could not generate Spacefile:\n    %s", err)
		}

		fmt.Printf("Spacefile generated at %s.\n", color.YellowString(filePath))
		fmt.Printf("Edit your Spacefile at will and use the %s command to actually create the new Space\n", color.YellowString("spacectl apply"))

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesInitCmd)

	spacesInitCmd.Flags().BoolVar(&spaceInitForce, "force", false, "Override existing Spacefile without asking")
	spacesInitCmd.Flags().StringVarP(&spaceInitName, "name", "n", "", "Name of the new Space")
	spacesInitCmd.Flags().StringVarP(&spaceInitLabel, "dns-label", "l", "", "DNS label of the new Space. Must be unique within the team.")

	spacesInitCmd.MarkFlagRequired("name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spacesInitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spacesInitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
