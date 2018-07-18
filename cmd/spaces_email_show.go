package cmd

import (
	"fmt"

	"os"

	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var spacesEmailShowFlags struct {
	SpaceID     string
	MessageID   string
	WithHeaders bool
	HTML        bool
}

var spacesEmailShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show an outgoing email",
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spacesEmailShowFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		emails, err := api.Spaces().ListCaughtEmails(space.ID)
		if err != nil {
			return err
		}

		email := emails.ByID(spacesEmailShowFlags.MessageID)
		if email == nil {
			return fmt.Errorf("no email with the ID '%s' was found", spacesEmailShowFlags.MessageID)
		}

		vh := view.CaughtEmailSingleView{
			CaughtEmail: *email,
			WithHeaders: spacesEmailShowFlags.WithHeaders,
			AsHTML:      spacesEmailShowFlags.HTML,
		}
		vh.Render(os.Stdout)

		return nil
	},
}

func init() {
	spaceEmailCmd.AddCommand(spacesEmailShowCmd)

	spacesEmailShowCmd.Flags().StringVarP(&spacesEmailShowFlags.SpaceID, "space", "s", "", "Space ID or name")
	spacesEmailShowCmd.Flags().StringVarP(&spacesEmailShowFlags.MessageID, "email", "m", "", "Message ID")
	spacesEmailShowCmd.Flags().BoolVar(&spacesEmailShowFlags.WithHeaders, "headers", false, "Display message headers")
	spacesEmailShowCmd.Flags().BoolVar(&spacesEmailShowFlags.HTML, "html", false, "Show HTML part if present")
}
