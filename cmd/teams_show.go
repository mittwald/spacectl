package cmd

import (
	"errors"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"os"
)

var teamShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details regarding a specific team",
	Long:  `Show details regarding a specific team`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			RootCmd.SilenceUsage = false
			return errors.New("Missing argument: Team ID or DNS label")
		}

		team, err := api.Teams().Get(args[0])
		if err != nil {
			return err
		}

		members, err := api.Teams().ListMembers(team.ID)
		if err != nil {
			return err
		}

		view := view.TabularTeamDetailView{IncludeMembers: true}
		view.TeamDetail(team, members, os.Stdout)

		return nil
	},
}

func init() {
	teamsCmd.AddCommand(teamShowCmd)
}
