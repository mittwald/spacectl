package view

import (
	"io"
	"fmt"
	"github.com/gosuri/uitable"
	"time"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/client/backups"
	"github.com/mittwald/spacectl/client/spaces"
)

type BackupView interface {
	List(backupList []backups.Backup, out io.Writer)
	Detail(backup *backups.Backup, recoveries []backups.Recovery, space *spaces.Space, out io.Writer)
}

type TabularBackupView struct {}

func (t TabularBackupView) List(backupList []backups.Backup, stage string, out io.Writer) {
	if len(backupList) == 0 {
		fmt.Fprintln(out, "No backups found.")
		return
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("ID", "STAGE", "STATUS", "KEEP", "DESCRIPTION", "SOFTWARE", "CREATED")

	for _, backup := range backupList {
		since := helper.HumanReadableDateDiff(time.Now(), backup.StartedAt)

		keep := "no"
		if backup.Keep {
			keep = "yes"
		}

		if stage == "" && backup.Stage != nil {
			stage = backup.Stage.Name
		}

		table.AddRow(
			backup.ID,
			stage,
			backup.Status,
			keep,
			backup.Description,
			fmt.Sprintf("%s %s", backup.Software.ID, backup.Version.Number),
			since+" ago",
		)
	}

	fmt.Fprintln(out, table)
}

func (t TabularBackupView) Detail(backup *backups.Backup, recoveries []backups.Recovery, space *spaces.Space, out io.Writer) {
	fmt.Fprintln(out, "GENERAL INFO")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	sinceStarted := helper.HumanReadableDateDiff(time.Now(), backup.StartedAt)

	table.AddRow("  ID:", backup.ID)
	table.AddRow("  Started:", sinceStarted + " ago (" + backup.StartedAt.Format(time.RFC1123) + ")")

	if !backup.CompletedAt.IsZero() {
		sinceCompleted := helper.HumanReadableDateDiff(time.Now(), backup.CompletedAt)

		table.AddRow("  Completed:", sinceCompleted + " ago (" + backup.CompletedAt.Format(time.RFC1123) + ")")
	} else {
		table.AddRow("  Completed:", "<pending>")
	}

	table.AddRow("  Status:", backup.Status)
	table.AddRow("  Description:", backup.Description)
	table.AddRow("  Application:", fmt.Sprintf("%s %s", backup.Software.ID, backup.Version.Number))

	if space != nil {
		table.AddRow("  Space:")
		table.AddRow("    ID:", space.ID)
		table.AddRow("    Name: ", fmt.Sprintf("%s (%s)", space.Name.HumanReadableName, space.Name.DNSName))
		table.AddRow("  Owning team:")
		table.AddRow("    ID:", space.Team.ID)
		table.AddRow("    Name:", space.Team.Name)
		table.AddRow("    DNS label:", space.Team.DNSLabel)
	} else {
		table.AddRow("  Space: <unknown>")
		table.AddRow("  Owning team: <unknown>")
	}

	fmt.Fprintln(out, table)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "RECOVERIES")

	recoveryTable := uitable.New()
	recoveryTable.Wrap = true
	recoveryTable.AddRow("  ID", "STATUS", "STARTED", "COMPLETED")

	format := "_2 Jan 2006 15:04 MST"

	for _, s := range recoveries {
		sinceStarted := helper.HumanReadableDateDiff(time.Now(), s.StartedAt)
		started := fmt.Sprintf("%s ago (%s)", sinceStarted, s.StartedAt.Format(format))
		completed := "<pending>"


		if !s.CompletedAt.IsZero() {
			sinceCompleted := helper.HumanReadableDateDiff(time.Now(), s.CompletedAt)
			completed = fmt.Sprintf("%s ago (%s)", sinceCompleted, s.CompletedAt.Format(format))
		}

		recoveryTable.AddRow(
			"  " + s.ID,
			s.Status,
			started,
			completed,
		)
	}

	fmt.Fprintln(out, recoveryTable)
}