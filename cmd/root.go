package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/mittwald/spacectl/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
	"strings"
	"github.com/fatih/color"
)

var cfgFile string
var apiServer string
var nonInteractive bool
var api client.SpacesClient
var verbose bool
var logger *log.Logger = log.New(ioutil.Discard, "spacectl", 0)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "spacectl",
	Short: "SPACES command line utility",
	Long:  `spacectl enables you to manage your SPACES hosting enviroment from the command line.`,
	SilenceErrors: true,
	SilenceUsage: true,

}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(`
     *           _______________________________
      /\        /                               \       *
     (  )      <  Houston, we have a problem... |
 *   (  )       \_______________________________/
    /|/\|\         *                                 *
   /_||||_\                 *
`)

		color.Red("")
		color.Red(err.Error())
		color.Red("")

		color.Blue("NEED HELP?")
		color.Blue("  * For support, contact support@spaces.de.")
		color.Blue("  * Report bugs with this tool at https://github.com/mittwald/spacectl")
		color.Blue("")

		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spaces/spaceconfig.yaml)")
	RootCmd.PersistentFlags().StringVar(&apiServer, "api-server", "https://api.dev.spaces.de", "API endpoint to connect to")
	RootCmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "Disable interactive prompts")
	RootCmd.PersistentFlags().String("token-file", "~/.spaces/token", "The file in which to store the authentication token")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Be more chatty")
	RootCmd.PersistentFlags().StringP("team", "t", "", "A team ID or name")

	viper.BindPFlag("apiServer", RootCmd.PersistentFlags().Lookup("api-server"))
	viper.BindPFlag("tokenFile", RootCmd.PersistentFlags().Lookup("token-file"))
	viper.BindPFlag("teamID", RootCmd.PersistentFlags().Lookup("team"))

	viper.BindEnv("apiServer", "SPACES_API_SERVER")
	viper.BindEnv("token", "SPACES_API_TOKEN")
	viper.BindEnv("teamID", "SPACES_TEAM_ID")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("spaceconfig")   // name of config file (without extension)
	viper.AddConfigPath("$HOME/.spaces") // adding home directory as first search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if token := viper.GetString("token"); token == "" {
		tokenFile := viper.GetString("tokenFile")

		if tokenFile[:2] == "~/" {
			usr, _ := user.Current()
			dir := usr.HomeDir
			tokenFile = filepath.Join(dir, tokenFile[2:])
		}

		c, err := ioutil.ReadFile(tokenFile)
		if err != nil && !os.IsNotExist(err) {
			panic(fmt.Errorf("could not read token file %s: %s", tokenFile, err))
		}

		viper.Set("token", strings.TrimSpace(string(c)))
	}

	clientConfig := client.SpacesClientConfig{
		Token: viper.GetString("token"),
		APIServer: viper.GetString("apiServer"),
		Logger: nil,
	}

	if verbose {
		clientConfig.Logger = log.New(os.Stderr, "spaces-client: ", 0)
		logger = log.New(os.Stderr, "spacectl: ", 0)
	}

	if s, err := client.NewSpacesClient(clientConfig); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	} else {
		api = s
	}
}
