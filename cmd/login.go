package cmd

import (
	"fmt"
	"github.com/mittwald/spacectl/service/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticates against the SPACES API",
	Long: `This command authenticates you against the SPACES API.

After logging in, this command will store an API token in ~/.spaces/token. Keep this file secret!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		service := auth.OAuthAuthenticationService{
			AuthServerURL: viper.GetString("authServer"),
		}

		result, err := service.Authenticate()
		if err != nil {
			return err
		}

		fmt.Println("successfully authenticated.")

		tokenFile := viper.GetString("tokenFile")

		if tokenFile[:2] == "~/" {
			usr, _ := user.Current()
			dir := usr.HomeDir
			tokenFile = filepath.Join(dir, tokenFile[2:])
		}

		tokenFileDir := filepath.Dir(tokenFile)

		_, err = os.Stat(tokenFileDir)
		if os.IsNotExist(err) {
			fmt.Printf("creating directory %s\n", tokenFileDir)
			os.MkdirAll(tokenFileDir, 0700)
		}

		ioutil.WriteFile(tokenFile, []byte(result.Token), 0600)
		fmt.Printf("token written to file %s\n", tokenFile)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().String("auth-server", "https://signup.spaces.de", "The URL of the SPACES authentication server")

	viper.BindPFlag("authServer", loginCmd.Flags().Lookup("auth-server"))
	viper.BindEnv("authServer", "SPACES_AUTH_SERVER")
}
