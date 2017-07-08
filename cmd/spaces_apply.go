package cmd

import (
	"github.com/mittwald/spacectl/spacefile"
	"github.com/spf13/cobra"
	"fmt"
)

var spacesApplySpacefile string

// applyCmd represents the apply command
var spacesApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies a Space configuration",
	Long: `This command reconciles a space definition from a Spacefile with the Spaces API.

CAUTION: This command can be potentially destructive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Printf("Using Spacefile at %s\n", spacesApplySpacefile)

		file, err := spacefile.ParseSpacefile(spacesApplySpacefile)
		if err != nil {
			return err
		}

		spc := file.Spaces[0]
		decl, err := spc.ToSpaceDeclaration()
		if err != nil {
			return err
		}

		declaredSpace, err := api.Spaces().Declare(spc.TeamID, decl)
		if err != nil {
			return err
		}

		fmt.Println(declaredSpace)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesApplyCmd)
	RootCmd.AddCommand(spacesApplyCmd)

	spacesApplyCmd.Flags().StringVarP(&spacesApplySpacefile, "spacefile", "f", "./" + spacefile.DefaultFilename, "Filename of Spacefile to apply")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
