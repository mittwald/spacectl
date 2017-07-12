package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mittwald/spacectl/spacefile"
)

var spacesValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a Spacefile for syntactical and semantical correctness",
	Long: `This command can be used to validate your Spacefile for syntactical and semantical correctness.

For example, it could be used before running "spacectl apply" (although "spacectl apply" will also
validate your Spacefile), or to check your Spacefile within your CI process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := spacefile.ParseSpacefile(spaceFile)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesValidateCmd)
}
