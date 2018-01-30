package cmd

import (
	"fmt"

	"errors"
	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/mittwald/spacectl/spacefile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"
)

var spaceInitFlags struct {
	Force    bool
	Name     string
	Label    string
	Software string
}

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

		if spaceInitFlags.Name == "" {
			mErr = multierror.Append(mErr, errors.New("must provide name (--name or -n)"))
		}

		if mErr != nil {
			return mErr
		}

		if spaceInitFlags.Label == "" {
			spaceInitFlags.Label = regexp.MustCompile("[^a-z0-9-]").ReplaceAllString(strings.ToLower(spaceInitFlags.Name), "-")
			fmt.Printf("Using %s as auto-generated DNS label\n", color.YellowString(spaceInitFlags.Label))
		}

		filePath := "./" + spacefile.DefaultFilename

		_, err := os.Stat(filePath)
		if !os.IsNotExist(err) {
			if !spaceInitFlags.Force {
				RootCmd.SilenceUsage = false
				return fmt.Errorf(`The file '%s' already exists in the current directory.
Use the --force flag (or -f) to overwrite it.`, filePath)
			} else {
				fmt.Printf("Overwriting existing Spacefile at %s\n", color.YellowString(filePath))
			}
		}

		sw, err := api.Applications().Get(spaceInitFlags.Software)
		if err != nil {
			return fmt.Errorf("Could not load software '%s':\n    %s", spaceInitFlags.Software, err)
		}

		fh, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("Could not open '%s' for writing:\n    %s", filePath, err)
		}

		err = spacefile.Generate(teamID, spaceInitFlags.Name, spaceInitFlags.Label, sw, fh)
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

	spacesInitCmd.Flags().BoolVar(&spaceInitFlags.Force, "force", false, "Override existing Spacefile without asking")
	spacesInitCmd.Flags().StringVarP(&spaceInitFlags.Name, "name", "n", "", "Name of the new Space")
	spacesInitCmd.Flags().StringVarP(&spaceInitFlags.Label, "dns-label", "l", "", "DNS label of the new Space. Must be unique within the team.")
	spacesInitCmd.Flags().StringVarP(&spaceInitFlags.Software, "software", "s", "typo3", "Software to initialize this Space with")

	spacesInitCmd.MarkFlagRequired("name")
}
