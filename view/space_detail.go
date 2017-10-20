package view

import (
	"github.com/mittwald/spacectl/client/spaces"
	"io"
	"fmt"
	"github.com/gosuri/uitable"
	"strings"
	"time"
	"github.com/mittwald/spacectl/cmd/helper"
)

type SpaceDetailView interface {
	SpaceDetail(space *spaces.Space, out io.Writer)
}

type TabularSpaceDetailView struct {}

func (t TabularSpaceDetailView) SpaceDetail(space *spaces.Space, out io.Writer) {
	fmt.Fprintln(out, "GENERAL INFO")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	since := helper.HumanReadableDateDiff(time.Now(), space.CreatedAt)

	table.AddRow("  ID:", space.ID)
	table.AddRow("  Created:", since + " ago")
	table.AddRow("  Created At:", space.CreatedAt.String())
	table.AddRow("  Status:", space.Status)
	table.AddRow("  Name:")
	table.AddRow("    Human-readable:", space.Name.HumanReadableName)
	table.AddRow("    DNS label:", space.Name.DNSName)
	table.AddRow("  Owning team:")
	table.AddRow("    ID:", space.Team.ID)
	table.AddRow("    Name:", space.Team.Name)
	table.AddRow("    DNS label:", space.Team.DNSLabel)

	fmt.Fprintln(out, table)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "STAGES")

	stageTable := uitable.New()
	stageTable.Wrap = true
	stageTable.AddRow("  NAME", "APPLICATION", "VERSION SPEC", "ACTUAL VERSION", "RUNNING", "DNS NAMES")

	for _, s := range space.Stages {
		running := "no"

		if s.Running {
			running = "yes"
		}

		stageTable.AddRow(
			"  " + s.Name,
			s.Application.ID,
			s.VersionConstraint,
			s.Version.Number,
			running,
			strings.Join(s.DNSNames, "\n"),
		)
	}

	fmt.Fprintln(out, stageTable)
}