package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/view"
	"bytes"
	"github.com/mittwald/spacectl/cmd/helper"
)

var spaceDeleteForce bool

var spaceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Space",
	Long: `This commnd deletes a Space.

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

	spaceDeleteCmd.Flags().BoolVar(&spaceDeleteForce, "force", false, "Do not prompt for confirmation")
}
