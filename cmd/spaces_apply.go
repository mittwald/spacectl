package cmd

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/mittwald/spacectl/service/costestimator"
	"github.com/mittwald/spacectl/view/confirm"
	"os"

	"github.com/mittwald/spacectl/client/spaces"

	"github.com/mittwald/spacectl/spacefile"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var spaceApplyFlags struct {
	AcceptTOS   bool
	AcceptCosts bool
}

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

		if !spaceApplyFlags.AcceptCosts {
			plans, err := api.Payment().ListPlans()
			if err != nil {
				return err
			}

			storage, err := spc.StorageBytes()
			if err != nil {
				return err
			}

			estimator := costestimator.New(plans)
			estimatorParams := costestimator.Params{
				PlanID: spc.Payment.PlanID,
				Stages: len(spc.Stages),
				StagesOnDemand: spc.CountOnDemandStages(),
				Scaling: decl.PaymentLink.Preprovisionings.Scaling.Quantity,
				Storage: storage,
				BackupIntervalMinutes: decl.PaymentLink.Options.BackupIntervalMinutes,
			}

			estimation, err := estimator.Estimate(estimatorParams)
			if err != nil {
				return err
			}

			buf := bytes.Buffer{}

			v := view.CostEstimationView{Estimation: *estimation}
			v.Render(&buf)

			c := confirm.Confirmation{
				Title: "Caution",
				Message: "Creating this Space will cause monthly costs (use the --accept-costs flag to skip this prompt).",
				Body: &buf,
				Color: color.New(color.FgYellow),
			}

			confirmed, err := c.DoPrompt(color.Output)
			if err != nil {
				return err
			}
			if confirmed == false {
				return nil
			}

			fmt.Println("")
		}

		if !spaceApplyFlags.AcceptTOS {
			c := confirm.Confirmation{
				Title: "Caution",
				Message: "Please accept to the Terms of Service before continuing (use the --accept-tos flag to skip this prompt).\nYou can find the current terms of service here: https://s3.eu-central-1.amazonaws.com/static.spaces.de/public/ABG_Mittwald_Final_Nov_2018.pdf",
				Color: color.New(color.FgYellow),
			}

			confirmed, err := c.DoPrompt(color.Output)
			if err != nil {
				return err
			}
			if confirmed == false {
				return nil
			}
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

			// check definition for virtualhosts and declare them
			for _, vhostDecl := range stageDef.VirtualHosts {
				vhost := vhostDecl.ToDeclaration()

				_, err = api.Spaces().UpdateVirtualHost(declaredSpace.ID, stageDecl.Name, vhost)
				if err != nil {
					return err
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

	spacesApplyCmd.Flags().BoolVar(&spaceApplyFlags.AcceptCosts, "accept-costs", false, "Agree to all occuring costs without being prompted")
	spacesApplyCmd.Flags().BoolVar(&spaceApplyFlags.AcceptTOS, "accept-tos", false, "Agree to terms of Service without being promted")
}
