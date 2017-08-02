package cmd

import (
	"fmt"

	"bytes"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var spaceDeleteForce bool

var spaceDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "Delete a Space",
	Long: `This command deletes a Space.

CAUTION: This command is destructive. Once you have deleted a Space, you
will not get it back!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		if !spaceDeleteForce {
			buf := bytes.Buffer{}
			view.TabularSpaceDetailView{}.SpaceDetail(space, &buf)

			ok, _ := view.Confirm("Once this Space is deleted, you will NOT be able to get it back.", buf.String())
			if !ok {
				fmt.Println("Aborting Space deletion.")
				return nil
			}
		}

		err = api.Spaces().Delete(space.ID)
		if err != nil {
			return err
		}

		fmt.Println("Space deleted.")

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spaceDeleteCmd)

	spaceDeleteCmd.Flags().BoolVarP(&spaceDeleteForce, "yes", "y", false, "Do not prompt for confirmation")
}
