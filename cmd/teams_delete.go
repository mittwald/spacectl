package cmd

import (
	"fmt"

	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"errors"
	"bytes"
)

var teamDeleteFlags struct {
	Force bool
}

var teamDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "Delete a team",
	Long: `This command deletes a team.

CAUTION: This command is destructive. Once you have deleted a team, you
will not get it back!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if len(args) < 1 {
			RootCmd.SilenceUsage = false
			return errors.New("Missing argument: Team ID or DNS label")
		}

		team, err := api.Teams().Get(args[0])
		if err != nil {
			return err
		}

		if !teamDeleteFlags.Force {
			buf := bytes.Buffer{}
			view.TabularTeamDetailView{}.TeamDetail(team, nil, &buf)

			ok, _ := view.Confirm("Once this Team is deleted, you will NOT be able to get it back.", buf.String())
			if !ok {
				fmt.Println("Aborting team deletion.")
				return nil
			}
		}

		err = api.Teams().Delete(args[0])
		if err != nil {
			return err
		}

		fmt.Println("Team deleted.")

		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamDeleteCmd)

	teamDeleteCmd.Flags().BoolVarP(&teamDeleteFlags.Force, "yes", "y", false, "Do not prompt for confirmation")
}
