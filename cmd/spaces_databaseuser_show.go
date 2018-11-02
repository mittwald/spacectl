package cmd

import (
	"errors"
	"os"

	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var databaseuserShowFlags struct {
	Username string
}

var spacesDatabaseuserShowCmd = &cobra.Command{
	Use:     "show --space <space-id> --stage <stage> --username <username>",
	Aliases: []string{"ls"},
	Short:   "Show database user",
	Long:    `Show a single database user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &databaseUserFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		stage, err := getStage(args, databaseUserFlags.StageName)
		if err != nil {
			return err
		}

		if databaseuserShowFlags.Username == "" {
			return errors.New("username expected: --username <username>")
		}

		user, err := api.Spaces().GetDatabaseUser(space.ID, stage, databaseuserShowFlags.Username)
		if err != nil {
			return err
		}

		v := view.TabularDatabaseUserView{}
		v.List([]spaces.DatabaseUser{*user}, os.Stdout)

		return nil
	},
}

func init() {
	spacesDatabaseuserCmd.AddCommand(spacesDatabaseuserShowCmd)

	f := spacesDatabaseuserShowCmd.PersistentFlags()
	f.StringVarP(&databaseuserShowFlags.Username, "username", "u", "", "Username of the new user")
}
