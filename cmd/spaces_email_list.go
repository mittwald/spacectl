package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var spacesEmailListFlags struct {
	SpaceID string
}

// listCmd represents the list command
var spacesEmailListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List outgoing emails",
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spacesEmailListFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		emails, err := api.Spaces().ListCaughtEmails(space.ID)
		if err != nil {
			return err
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("DATE", "SENDER", "RECIPIENT", "SUBJECT")

		for _, email := range emails {
			vh := view.CaughtEmailView{CaughtEmail: email}

			table.AddRow(
				vh.RenderDate(),
				vh.RenderSender(),
				vh.RenderRecipients(2),
				vh.Mail.Subject,
			)
		}

		fmt.Println(table)

		return nil
	},
}

func init() {
	spaceEmailCmd.AddCommand(spacesEmailListCmd)

	spacesEmailListCmd.Flags().StringVarP(&spacesEmailListFlags.SpaceID, "space", "s", "", "Space ID or name")
}
