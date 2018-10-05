package cmd

import (
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"os"
)

var spacesDatabaseuserListCmd = &cobra.Command{
	Use:     "list --space <space-id> --stage <stage>",
	Aliases: []string{"ls"},
	Short:   "List database users",
	Long:    `Lists database users`,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := api.Spaces().ListDatabaseUsersByStage(databaseUserFlags.SpaceID, databaseUserFlags.StageName)
		if err != nil {
			return err
		}

		v := view.TabularDatabaseUserView{}
		v.List(users, os.Stdout)

		return nil
	},
}

func init() {
	spacesDatabaseuserCmd.AddCommand(spacesDatabaseuserListCmd)
}
