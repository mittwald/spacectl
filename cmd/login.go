// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"errors"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"fmt"
	"github.com/mittwald/spacectl/service/auth"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login -u <username> [-p <password>]",
	Short: "Authenticates against the SPACES API",
	Long: `This command authenticates you against the SPACES API.

After logging in, this command will store an API token in ~/.spaces/token. Keep this file secret!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := false
		username := cmd.Flag("username").Value.String()

		if username == "" {
			return errors.New("No --username specified")
		}

		for {
			password := cmd.Flag("password").Value.String()
			if password == "" {
				if nonInteractive {
					return errors.New("No --password specified")
				}

				prompt = true

				fmt.Print("Enter password: ")
				passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))

				fmt.Println("")

				if err != nil {
					return err
				}

				password = string(passwordBytes)
			}

			service := auth.AuthenticationService{
				AuthServerURL: cmd.Flag("auth-server").Value.String(),
			}

			_, err := service.Authenticate(username, password)
			if err != nil {
				switch tErr := err.(type) {
				case auth.InvalidCredentialsErr:
					if nonInteractive || !prompt {
						return tErr
					}

					fmt.Println("invalid credentials. try again")
				default:
					return tErr
				}
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "The username with which to connect")
	loginCmd.Flags().StringP("password", "p", "", "The password with which to connect")
	loginCmd.Flags().String("auth-server", "https://signup.dev.spaces.de", "The URL of the SPACES authentication server")
}
