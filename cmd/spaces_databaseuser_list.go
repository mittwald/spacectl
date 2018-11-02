package cmd

import (
	"errors"
	"os"

	"github.com/hashicorp/go-multierror"

	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var spacesDatabaseuserListCmd = &cobra.Command{
	Use:     "list --space <space-id> --stage <stage>",
	Aliases: []string{"ls"},
	Short:   "List database users",
	Long:    `Lists database users`,
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &databaseUserFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		var stage string
		if databaseUserFlags.StageName != "" {
			stage = databaseUserFlags.StageName
		} else if len(args) >= 2 {
			stage = args[1]
		} else {
			return multierror.Append(
				errors.New("stage name expected. use"),
				errors.New("-e [STAGE_NAME] or"),
				errors.New("list [SPACE_ID] [STAGE_ID]"),
			)
		}

		users, err := api.Spaces().ListDatabaseUsersByStage(space.ID, stage)
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
