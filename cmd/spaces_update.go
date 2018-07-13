package cmd

import (
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"os"
)

var spaceUpdateFlags struct {
	SpaceIDOrName string
	SourceStage   string
	TargetStage   string
	Version       string
}

var spaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Triggers a software update ",
	Long:  `This command starts a software update in the specified stage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spaceOpenFlags.SpaceIDOrName, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		update, err := api.Spaces().UpdateApplication(
			space.ID,
			spaceUpdateFlags.SourceStage,
			spaceUpdateFlags.TargetStage,
			spaceUpdateFlags.Version,
		)

		if err != nil {
			return err
		}

		v := view.TabularSpaceApplicationUpdateView{}
		v.SpaceApplicationUpdate(space, update, os.Stdout)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spaceUpdateCmd)

	f := spaceUpdateCmd.Flags()
	f.StringVarP(&spaceUpdateFlags.SourceStage, "stage", "e", "production", "The stage to update")
	f.StringVar(&spaceUpdateFlags.TargetStage, "target-stage", "", "The stage to execute the update in. Default to the value of --stage if omitted.")
	f.StringVarP(&spaceUpdateFlags.SpaceIDOrName, "space", "s", "", "The space in which to run the update")
	f.StringVar(&spaceUpdateFlags.Version, "version", "", "The version constraint to update to")
}
