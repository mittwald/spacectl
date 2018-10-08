package cmd

import (
	"os"

	"github.com/mittwald/spacectl/client/spaces"

	"github.com/mittwald/spacectl/spacefile"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var spacesApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies a Space configuration",
	Long: `This command reconciles a space definition from a Spacefile with the Spaces API.

CAUTION: This command can be potentially destructive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Printf("Using Spacefile at %s\n", spaceFile)

		file, err := spacefile.ParseSpacefile(spaceFile, false)
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

		// get declared stages and find their definition
		for _, stageDecl := range declaredSpace.Stages {
			stageDef := spc.GetStageByName(stageDecl.Name)
			if stageDef == nil {
				continue
			}

			if len(stageDef.VirtualHosts) == 0 {
				// check definition for virtualhosts and declare them
				for _, vhostDecl := range stageDef.VirtualHosts {
					vhost := vhostDecl.ToDeclaration()

					_, err = api.Spaces().UpdateVirtualHost(declaredSpace.ID, stageDecl.Name, vhost)
					if err != nil {
						return err
					}
				}
			}

			if stageDef.Protection != "" {
				_, err = api.Spaces().CreateStageProtection(declaredSpace.ID, stageDecl.Name, spaces.StageProtection{ProtectionType: stageDef.Protection})
			} else {
				err = api.Spaces().DeleteStageProtection(declaredSpace.ID, stageDecl.Name)
			}

			if err != nil {
				return err
			}
		}

		updates, err := api.Spaces().ListApplicationUpdatesBySpace(declaredSpace.ID)
		if err != nil {
			return err
		}

		payment, err := api.Spaces().GetPaymentLink(declaredSpace.ID)
		if err != nil {
			return err
		}

		v := view.TabularSpaceDetailView{}
		v.SpaceDetail(declaredSpace, updates, payment, os.Stdout)

		return nil
	},
}

func init() {
	spacesCmd.AddCommand(spacesApplyCmd)
}
