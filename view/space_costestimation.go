package view

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/service/costestimator"
	"golang.org/x/text/currency"
	"io"
	"strings"
)

type CostEstimationView struct {
	costestimator.Estimation
}

func (v *CostEstimationView) Render(out io.Writer) {
	y := color.New(color.FgYellow)

	fmt.Fprintln(out, "ESTIMATED MONTHLY CHARGES")

	table := uitable.New()
	table.MaxColWidth = 120
	table.Wrap = true
	table.RightAlign(3)

	table.AddRow("  #", "ITEM", "QUANTITY", "MONTHLY COST")

	var longest [3]int

	for i := range v.LineItems {
		table.AddRow(fmt.Sprintf("  %d", i + 1), v.LineItems[i].Subject, v.LineItems[i].Quantity.String(), currency.ISO(v.LineItems[i].MonthlyCost))

		longest[0] = max(longest[0], len(v.LineItems[i].Subject))
		longest[1] = max(longest[1], len(v.LineItems[i].Quantity.String()))
		longest[2] = max(longest[2], len(fmt.Sprintf("%s", currency.ISO(v.LineItems[i].MonthlyCost))))
	}

	table.AddRow("", strings.Repeat("-", longest[0]), strings.Repeat("-", longest[1]), strings.Repeat("-", longest[2]))
	table.AddRow("", "TOTAL MONTHLY COST", "", currency.ISO(v.MonthlyTotalCost))

	fmt.Fprintln(out, table)
	fmt.Fprintln(out, "")

	y.Fprintln(out, "ON-DEMAND CHARGES")
	y.Fprintln(out, "  If you exceed your allocated resources, additional charges may apply:")
	y.Fprintf(out, "    * per additional GiB of used storage: %s per month and GiB\n", currency.ISO(v.OnDemandCharges.Storage.Amount))
	y.Fprintf(out, "    * per additional development stage:   %s per month and stage\n", currency.ISO(v.OnDemandCharges.Stages.Amount))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}