package cmd

import (
	"fmt"

	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/spacefile"
	"github.com/rodaine/hclencoder"
	"github.com/spf13/cobra"
)

var spacesExportFlags struct {
	SpaceID string
}

var spacesExportCmd = &cobra.Command{
	Use:   "export <space>",
	Short: "export a existing space",
	Long: `this command exports an existing space to a spacefile.

	This spacefile can be edited and applied again using spacectl space apply.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spacesExportFlags.SpaceID, api)
		if err != nil {
			return err
		}

		spaceDef := spacefile.FromSpace(space)

		hcl, err := hclencoder.Encode(&spaceDef)
		if err != nil {
			return err
		}

		fmt.Println(string(hcl))

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesExportCmd)
	spacesExportCmd.Flags().StringVarP(&spacesExportFlags.SpaceID, "space", "s", "", "Space ID or name")
}
