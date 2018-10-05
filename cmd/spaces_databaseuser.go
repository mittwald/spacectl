package cmd

import "github.com/spf13/cobra"

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

func init() {
	spacesCmd.AddCommand(spacesDatabaseuserCmd)

	f := spacesDatabaseuserCmd.PersistentFlags()
	f.StringVarP(&databaseUserFlags.SpaceID, "space", "s", "", "Space ID")
	f.StringVarP(&databaseUserFlags.StageName, "stage", "e", "", "Stage name")
}
