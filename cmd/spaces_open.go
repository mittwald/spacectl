package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/spacefile"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/client/spaces"
	"strings"
)

var spaceOpenStage string = "production"

// openCmd represents the open command
var spaceOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Opens the Space in your browser",
	Long: `This command opens a Space in your browser.

This command will respect a ` + spacefile.DefaultFilename + ` file in your current directory.
Alternatively, use the -t and -s flags to specify team and space ID/name.

By default, this command will open the Space's first defined stage
(typically, "production"). To change this, supply the --stage or -e flag.'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		var stage *spaces.Stage
		var existingStageNames []string = make([]string, len(space.Stages))

		for i := range space.Stages {
			if space.Stages[i].Name == spaceOpenStage {
				stage = &space.Stages[i]
			}
			existingStageNames[i] = space.Stages[i].Name
		}

		if stage == nil {
			RootCmd.SilenceUsage = false
			return fmt.Errorf(
				"The Space '%s' does not have a stage '%s'. Existing stages are: '%s'.",
				space.ID,
				spaceOpenStage,
				strings.Join(existingStageNames, "', '"),
			)
		}

		if len(stage.DNSNames) == 0 {
			return fmt.Errorf("We do not have an URL for this stage, yet. Please try again in a few seconds.")
		}

		url := "http://" + stage.DNSNames[0]
		fmt.Printf("Opening %s\n", url)

		helper.OpenURL(url)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spaceOpenCmd)

	spaceOpenCmd.Flags().StringVarP(&spaceOpenStage, "stage", "e", "production", "The stage to open")
}
