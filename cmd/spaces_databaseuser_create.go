package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/hashicorp/go-multierror"

	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var databaseUserCreateFlags struct {
	Type     string
	Username string
	External string
}

var spacesDatabaseuserCreateCmd = &cobra.Command{
	Use:     "create --space <space-id> --stage <stage> [--external <ip_schema>] --type <type> --username <username>",
	Aliases: []string{"add"},
	Short:   "Create a database user",
	Long:    `Creates a new user in the given database`,
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

		if databaseUserCreateFlags.Type == "" || databaseUserCreateFlags.Username == "" {
			return multierror.Append(
				errors.New("type and username are required"),
				errors.New("--type <type> --username <username>"),
			)
		}

		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()
		password := strings.TrimSpace(string(bytePassword))

		dbUser := spaces.DatabaseUserInput{
			UserSuffix: databaseUserCreateFlags.Username,
			Password:   password,
			Type:       databaseUserCreateFlags.Type,
			External:   databaseUserCreateFlags.External,
		}

		user, err := api.Spaces().CreateDatabaseUser(space.ID, stage, dbUser)
		if err != nil {
			return err
		}

		v := view.TabularDatabaseUserView{}
		v.List([]spaces.DatabaseUser{*user}, os.Stdout)

		return nil
	},
}

func init() {
	spacesDatabaseuserCmd.AddCommand(spacesDatabaseuserCreateCmd)

	f := spacesDatabaseuserCreateCmd.PersistentFlags()
	f.StringVarP(&databaseUserCreateFlags.Type, "type", "d", "mysql", "Creates the user in this database (default: mysql)")
	f.StringVarP(&databaseUserCreateFlags.Username, "username", "u", "", "Username of the new user")
	f.StringVarP(&databaseUserCreateFlags.External, "external", "", "", "Authorized external IPs. If omitted, the user is only able to login from the application pod. '%' is a wildcard: e.g. 192.168.%")
}
