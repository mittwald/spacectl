package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/spacefile"
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/spf13/viper"
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/mittwald/spacectl/view"
	"bytes"
)

var spaceDeleteForce bool
var spaceDeleteFile string

var spaceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Space",
	Long: `This commnd deletes a Space.

CAUTION: This command is destructive. Once you have deleted a Space, you
will not get it back!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := func() (*spaces.Space, error) {
			if len(args) >= 1 {
				teamID := viper.GetString("teamID")
				if teamID != "" {
					return api.Spaces().GetByTeamAndName(teamID, args[0])
				} else {
					return api.Spaces().GetByID(args[0])
				}
			}

			f, err := spacefile.ParseSpacefile(spaceDeleteFile)
			if err == nil {
				spaceDef := f.Spaces[0]

				space, err := api.Spaces().GetByTeamAndName(spaceDef.TeamID, spaceDef.DNSLabel)
				if err != nil {
					return nil, err
				}

				return space, nil
			}

			if _, ok := err.(spacefile.ErrSpacefileNotFound); ok {
				RootCmd.SilenceUsage = false
				err := multierror.Append(nil,
					fmt.Errorf("No spacefile found at %s", spaceDeleteFile),
					errors.New("Missing team ID (--team, -t or $SPACES_TEAM_ID)"),
				)
				return nil, err
			}

			return nil, err
		}()

		if err != nil {
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
	spaceDeleteCmd.Flags().StringVarP(&spaceDeleteFile, "spacefile", "f", "./" + spacefile.DefaultFilename, "Delete Space defined in Spacefile")
}
