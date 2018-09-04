package cmd

import (
	"fmt"

	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/spacefile"
	"github.com/rodaine/hclencoder"
	"github.com/spf13/cobra"
)

var spacesDescribeFlags struct {
	SpaceID string
}

var spacesDescribeCmd = &cobra.Command{
	Use:   "describe <space>",
	Short: "Describe a existing space",
	Long: `This command describes an existing space.

	Outputs the spacefile for an existing space. This spacefile can be applied
	again using spacectl space apply.
	A space generated from the web frontend can be modified on the cli and updated
	with the apply command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spacesDescribeFlags.SpaceID, api)
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
	spacesCmd.AddCommand(spacesDescribeCmd)
	spacesDescribeCmd.Flags().StringVarP(&spacesDescribeFlags.SpaceID, "space", "s", "", "Space ID or name")
}
