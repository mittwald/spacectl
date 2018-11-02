package cmd

import (
	"errors"

	"github.com/hashicorp/go-multierror"

	"github.com/spf13/cobra"
)

var databaseUserFlags struct {
	SpaceID   string
	StageName string
}

var spacesDatabaseuserCmd = &cobra.Command{
	Use:     "databaseusers",
	Aliases: []string{"databaseuser", "dbusers", "dbuser", "dbu"},
	Short:   "Manage database users",
	Long:    `This command allows you to manage database users`,
}

func getStage(args []string, option string) (string, error) {
	var stage string
	if option != "" {
		stage = option
	} else if len(args) >= 2 {
		stage = args[1]
	} else {
		return "", multierror.Append(
			errors.New("stage name expected. use"),
			errors.New("-e [STAGE_NAME] or"),
			errors.New("list [SPACE_ID] [STAGE_ID]"),
		)
	}

	return stage, nil
}

func init() {
	spacesCmd.AddCommand(spacesDatabaseuserCmd)

	f := spacesDatabaseuserCmd.PersistentFlags()
	f.StringVarP(&databaseUserFlags.SpaceID, "space", "s", "", "Space ID")
	f.StringVarP(&databaseUserFlags.StageName, "stage", "e", "", "Stage name")
}
