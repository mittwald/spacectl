package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

var plansListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List plans",
	Long:    `Lists all available plans that you can use.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		plans, err := api.Payment().ListPlans()
		if err != nil {
			return err
		}

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("ID", "NAME", "PRICE")

		for _, plan := range plans {
			table.AddRow(plan.ID, plan.Name, fmt.Sprintf("%6.2f %s/%s", plan.BasePrice.Value, plan.BasePrice.Currency.Unit(), plan.BasePrice.Unit.String()))
		}

		fmt.Println(table)

		return nil
	},
}

func init() {
	plansCmd.AddCommand(plansListCmd)
}
