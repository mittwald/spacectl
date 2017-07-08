package cmd

import (
	"github.com/spf13/cobra"
	"github.com/hashicorp/go-multierror"
	"errors"
	"github.com/mittwald/spacectl/view"
	"os"
)

// spacesShowCmd represents the show command
var spacesShowCmd = &cobra.Command{
	Use:   "show -t <team> <space-name>",
	Short: "Show details regarding a specific space",
	Long: "Show details regarding a specific space",
	RunE: func(cmd *cobra.Command, args []string) error {
		var mErr *multierror.Error
		//teamID := viper.GetString("teamID")

		if len(args) == 0 {
			mErr = multierror.Append(mErr, errors.New("missing Space identifier (either ID or DNS label)"))
		}

		if mErr.ErrorOrNil() != nil {
			RootCmd.SilenceUsage = false
			return mErr
		}

		space, err := api.Spaces().GetByID(args[0])
		if err != nil {
			return err
		}

		v := view.TabularSpaceDetailView{}
		v.SpaceDetail(space, os.Stdout)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesShowCmd)
}
