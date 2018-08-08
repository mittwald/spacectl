package cmd

import (
	"fmt"

	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/view"
	"github.com/spf13/cobra"
)

var spacesTopFlags struct {
	SpaceID string
	Scope   string
}

var spacesTopCmd = &cobra.Command{
	Use:   "top -t <team> -s <space-name>",
	Short: "Show usage statistics for a certain Space",
	RunE: func(cmd *cobra.Command, args []string) error {
		space, err := helper.GetSpaceFromContext(args, spaceFile, &spacesTopFlags.SpaceID, api)
		if err != nil {
			RootCmd.SilenceUsage = false
			return err
		}

		metrics, err := api.Spaces().GetComputeMetrics(space.ID, spacesTopFlags.Scope)
		if err != nil {
			return err
		}

		if len(metrics) == 0 {
			fmt.Println("no metrics found")
			return nil
		}

		v := view.SpaceMetricsView{}
		return v.Render(spacesTopFlags.Scope, func() (spaces.ComputeMetricPointList, error) {
			return api.Spaces().GetComputeMetrics(space.ID, spacesTopFlags.Scope)
		})
	},
}

func init() {
	spacesCmd.AddCommand(spacesTopCmd)
	spacesTopCmd.Flags().StringVarP(&spacesTopFlags.SpaceID, "space", "s", "", "Space ID or name")
	spacesTopCmd.Flags().StringVar(&spacesTopFlags.Scope, "scope", spaces.ScopeWeek, "Time scope (hour, today, week, month, year)")
}
