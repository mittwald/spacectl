package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/view"
	"os"
	"github.com/mittwald/spacectl/cmd/helper"
)

var spacesShowFlags struct {
	SpaceID string
}

var spacesShowCmd = &cobra.Command{
	Use:   "show -t <team> <space-name>",
	Short: "Show details regarding a specific space",
	Long: "Show details regarding a specific space",
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spacesShowFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		updates, err := api.Spaces().ListApplicationUpdatesBySpace(space.ID)
		if err != nil {
			return err
		}

		v := view.TabularSpaceDetailView{}
		v.SpaceDetail(space, updates, os.Stdout)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesShowCmd)
	spacesShowCmd.Flags().StringVarP(&spacesShowFlags.SpaceID, "space", "s", "", "Space ID or name")
}
