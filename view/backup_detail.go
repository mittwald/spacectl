package view

import (
	"io"
	"fmt"
	"github.com/gosuri/uitable"
	"time"
	"github.com/mittwald/spacectl/cmd/helper"
	"github.com/mittwald/spacectl/client/backups"
)

type BackupDetailView interface {
	BackupDetail(backup *backups.Backup, recoveries []backups.Recovery, out io.Writer)
}

type TabularBackupDetailView struct {}

func (t TabularBackupDetailView) BackupDetail(backup *backups.Backup, recoveries []backups.Recovery, out io.Writer) {
	fmt.Fprintln(out, "GENERAL INFO")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	sinceStarted := helper.HumanReadableDateDiff(time.Now(), backup.StartedAt)

	table.AddRow("  ID:", backup.ID)
	table.AddRow("  Started:", sinceStarted + " ago")
	table.AddRow("  Started At:", backup.StartedAt.String())

	if !backup.CompletedAt.IsZero() {
		sinceCompleted := helper.HumanReadableDateDiff(time.Now(), backup.CompletedAt)

		table.AddRow("  Completed:", sinceCompleted + " ago")
		table.AddRow("  Completed At:", backup.CompletedAt.String())
	} else {
		table.AddRow("  Completed:", "<pending>")
		table.AddRow("  Completed At:", "<pending>")
	}

	table.AddRow("  Status:", backup.Status)
	table.AddRow("  Description:", backup.Description)

	fmt.Fprintln(out, table)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "RECOVERIES")

	recoveryTable := uitable.New()
	recoveryTable.Wrap = true
	recoveryTable.AddRow("  ID", "STATUS", "STARTED", "COMPLETED")

	for _, s := range recoveries {
		sinceStarted := helper.HumanReadableDateDiff(time.Now(), backup.StartedAt)
		sinceCompleted := helper.HumanReadableDateDiff(time.Now(), backup.CompletedAt)

		if backup.CompletedAt.IsZero() {
			sinceCompleted = "<pending>"
		}

		recoveryTable.AddRow(
			"  " + s.ID,
			s.Status,
			sinceStarted,
			sinceCompleted,
		)
	}

	fmt.Fprintln(out, recoveryTable)
}