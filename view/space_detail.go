package view

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/cmd/helper"
)

type SpaceDetailView interface {
	SpaceDetail(space *spaces.Space, updates []spaces.ApplicationUpdate, out io.Writer)
}

type TabularSpaceDetailView struct{}

func (t TabularSpaceDetailView) SpaceDetail(space *spaces.Space, updates []spaces.ApplicationUpdate, paymentLink *spaces.SpacePaymentLink, out io.Writer) {
	fmt.Fprintln(out, "GENERAL INFO")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	since := helper.HumanReadableDateDiff(time.Now(), space.CreatedAt)

	table.AddRow("  ID:", space.ID)
	table.AddRow("  Created:", since+" ago")
	table.AddRow("  Created At:", space.CreatedAt.String())
	table.AddRow("  Name:")
	table.AddRow("    Human-readable:", space.Name.HumanReadableName)
	table.AddRow("    DNS label:", space.Name.DNSName)
	table.AddRow("  Owning team:")
	table.AddRow("    ID:", space.Team.ID)
	table.AddRow("    Name:", space.Team.Name)
	table.AddRow("    DNS label:", space.Team.DNSLabel)
	table.AddRow("  Payment settings:")

	if paymentLink == nil {
		table.AddRow("    (not configured yet)")
	} else {
		p := &paymentLink.Plan
		pr := &paymentLink.PaymentProfile
		table.AddRow("    Plan:", fmt.Sprintf("%s (%s): %.2f %s/%s", p.Name, p.ID, p.BasePrice.Value, p.BasePrice.Currency.Unit(), p.BasePrice.Unit.String()))
		table.AddRow("    Payment profile:")
		table.AddRow("      ID:", pr.ID)

		if pr.ContractPartner.Company != "" {
			table.AddRow("      Contract Partner:", pr.ContractPartner.Company+", "+pr.ContractPartner.FirstName+" "+pr.ContractPartner.LastName)
		} else {
			table.AddRow("      Contract Partner:", pr.ContractPartner.FirstName+" "+pr.ContractPartner.LastName)
		}
	}

	fmt.Fprintln(out, table)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "STAGES")

	stageTable := uitable.New()
	stageTable.Wrap = true
	stageTable.AddRow("  NAME", "APPLICATION", "VERSION SPEC", "ACTUAL VERSION", "STATUS", "DNS NAMES")

	for _, s := range space.Stages {
		stageTable.AddRow(
			"  "+s.Name,
			s.Application.Software.ID,
			s.VersionConstraint,
			s.Version.Number,
			s.Status,
			strings.Join(s.DNSNames, "\n"),
		)
	}

	fmt.Fprintln(out, stageTable)
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "APPLICATION UPDATES")

	if len(updates) == 0 {
		fmt.Fprintln(out, "  No application updates found")
	} else {
		updateTable := uitable.New()
		updateTable.Wrap = true
		updateTable.AddRow("  ID", "STARTED", "COMPLETED", "VERSION", "SOURCE STAGE", "TARGET STAGE", "PROGRESS")

		for _, u := range updates {
			started := helper.HumanReadableDateDiff(time.Now(), u.StartedAt) + " ago"
			completed := "(pending)"
			progress := "(pending)"

			if !u.CompletedAt.IsZero() {
				completed = helper.HumanReadableDateDiff(time.Now(), u.CompletedAt) + " ago"
			}

			if u.Progress.TotalSteps > 0 {
				progress = fmt.Sprintf("%d/%d (%s)", u.Progress.CurrentStep, u.Progress.TotalSteps, u.Progress.Status)
			}

			updateTable.AddRow(
				"  "+u.ID,
				started,
				completed,
				u.ExactVersion.Number+" ("+u.VersionConstraint+")",
				u.SourceStage.Name,
				u.TargetStage.Name,
				progress,
			)
		}

		fmt.Fprintln(out, updateTable)
	}
}
