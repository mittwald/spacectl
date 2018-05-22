package view

import (
	"github.com/mittwald/spacectl/client/spaces"
	"io"
	"fmt"
	"github.com/gosuri/uitable"
	"time"
	"github.com/mittwald/spacectl/cmd/helper"
)

type SpaceApplicationUpdateView interface {
	SpaceApplicationUpdate(space *spaces.Space, update *spaces.ApplicationUpdate, out io.Writer)
}

type TabularSpaceApplicationUpdateView struct {}

func (t TabularSpaceApplicationUpdateView) SpaceApplicationUpdate(space *spaces.Space, update *spaces.ApplicationUpdate, out io.Writer) {
	fmt.Fprintln(out, "GENERAL INFO")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	startedRelative := helper.HumanReadableDateDiff(time.Now(), update.StartedAt)
	completedRelative := "(pending)"
	completed := "(pending)"

	if !update.CompletedAt.IsZero() {
		completedRelative = helper.HumanReadableDateDiff(time.Now(), update.CompletedAt) + " ago"
		completed = update.CompletedAt.String()
	}

	table.AddRow("  ID:", space.ID)
	table.AddRow("  Started:", startedRelative+ " ago")
	table.AddRow("  Started At:", update.StartedAt.String())
	table.AddRow("  Completed:", completedRelative)
	table.AddRow("  Completed At:", completed)
	table.AddRow("  Status:", update.Status)
	table.AddRow("  Space:")
	table.AddRow("    Human-readable:", space.Name.HumanReadableName)
	table.AddRow("    DNS label:", space.Name.DNSName)
	table.AddRow("  Source stage:", update.SourceStage.Name)
	table.AddRow("  Target stage:", update.TargetStage.Name)

	if update.Progress.TotalSteps > 0 {
		table.AddRow("  Progress:")
		table.AddRow("    Step:", fmt.Sprintf("%d / %d", update.Progress.CurrentStep, update.Progress.TotalSteps))
		table.AddRow("    Status:", update.Progress.Status)
	}

	fmt.Fprintln(out, table)
}