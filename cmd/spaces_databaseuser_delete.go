package cmd

import (
	"errors"

	"github.com/hashicorp/go-multierror"

	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/spf13/cobra"
)

var databaseUserDeleteFlags struct {
	Type     string
	Username string
}

var spacesDatabaseuserDeleteCmd = &cobra.Command{
	Use:     "delete --space <space-id> --stage <stage> --type <type> --username <username>",
	Aliases: []string{"del", "d"},
	Short:   "Deletes a database user",
	Long:    `Deletes a new user`,
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

		if databaseUserDeleteFlags.Type == "" || databaseUserDeleteFlags.Username == "" {
			return multierror.Append(
				errors.New("type and username are required"),
				errors.New("--type <type> --username <username>"),
			)
		}

		err = api.Spaces().DeleteDatabaseUser(space.ID, stage, databaseUserDeleteFlags.Username, databaseUserDeleteFlags.Type)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	spacesDatabaseuserCmd.AddCommand(spacesDatabaseuserDeleteCmd)

	f := spacesDatabaseuserDeleteCmd.PersistentFlags()
	f.StringVarP(&databaseUserDeleteFlags.Type, "type", "d", "mysql", "Creates the user in this database (default: mysql)")
	f.StringVarP(&databaseUserDeleteFlags.Username, "username", "u", "", "Username of the new user")
}
