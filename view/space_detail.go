package view

import (
	"github.com/mittwald/spacectl/client/spaces"
	"io"
	"fmt"
	"github.com/gosuri/uitable"
	"strings"
	"time"
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

	fmt.Printf("%v\n", space)

	round := time.Second
	since := time.Now().Round(round).Sub(space.CreatedAt.Round(round)).String()

	table.AddRow("  ID:", space.ID)
	table.AddRow("  Created:", since + " ago")
	table.AddRow("  Created At:", space.CreatedAt.String())
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
	stageTable.AddRow("  NAME", "APPLICATION", "VERSION SPEC", "ACTUAL VERSION", "DNS NAMES")

	for _, s := range space.Stages {
		stageTable.AddRow(
			"  " + s.Name,
			s.Application.ID,
			s.VersionConstraint,
			s.Version.Number,
			strings.Join(s.DNSNames, "\n"),
		)
	}

	fmt.Fprintln(out, stageTable)
}