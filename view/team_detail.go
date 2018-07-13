package view

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/teams"
	"github.com/mittwald/spacectl/cmd/helper"
	"io"
	"strings"
	"time"
)

type TeamDetailView interface {
	TeamDetail(team *teams.Team, members []teams.Membership, out io.Writer)
}

type TabularTeamDetailView struct {
	IncludeMembers bool
}

func (t TabularTeamDetailView) TeamDetail(team *teams.Team, members []teams.Membership, out io.Writer) {
	fmt.Fprintln(out, "GENERAL INFO")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	since := helper.HumanReadableDateDiff(time.Now(), team.CreatedAt)

	table.AddRow("  ID:", team.ID)
	table.AddRow("  Created:", since+" ago")
	table.AddRow("  Created At:", team.CreatedAt.String())
	table.AddRow("  Name:")
	table.AddRow("    Human-readable:", team.Name)
	table.AddRow("    DNS label:", team.DNSName)

	fmt.Fprintln(out, table)

	if t.IncludeMembers {
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "MEMBERS")

		memberTable := uitable.New()
		memberTable.Wrap = true
		memberTable.AddRow("  ID", "ROLE", "NAME", "E-MAIL")

		for _, m := range members {
			name := strings.TrimSpace(fmt.Sprintf("%s %s", m.User.FirstName, m.User.LastName))
			memberTable.AddRow("  "+m.User.ID, m.Role, name, m.User.Email)
		}

		fmt.Fprintln(out, memberTable)
	}
}
