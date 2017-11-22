package cmd

import (
	"github.com/spf13/cobra"
	"errors"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"github.com/mittwald/spacectl/view"
	"os"
)

var keyImportFlags struct {
	Comment string
}

var keyImportCmd = &cobra.Command{
	Use:   "import <key-file>",
	Short: "Import an existing SSH public key",
	Long: `This command imports an existing SSH public key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			RootCmd.SilenceUsage = false
			return errors.New("Missing argument: Key file name")
		}

		keyFile := args[0]
		keyBytes, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return err
		}

		key, comment, _, _, err := ssh.ParseAuthorizedKey(keyBytes)
		if err != nil {
			return err
		}

		if keyImportFlags.Comment != "" {
			comment = keyImportFlags.Comment
		}

		createdKey, err := api.SSHKeys().Add(key.Marshal(), key.Type(), comment)
		if err != nil {
			return err
		}

		v := view.TabularKeyDetailView{}
		v.KeyDetail(createdKey, os.Stdout)

		return nil
	},
}

func init() {
	keysCmd.AddCommand(keyImportCmd)

	keyImportCmd.Flags().StringVarP(&keyImportFlags.Comment, "comment", "c", "", "Public key comment (if left out, the comment from the public key file will be used)")
}
