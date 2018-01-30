package view

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/software"
	"io"
)

type AppsView interface {
	List(softwareList []software.Software, out io.Writer)
}

type TabularAppsView struct {}

func (t TabularAppsView) List(softwareList []software.Software, out io.Writer) {
	if len(softwareList) == 0 {
		fmt.Fprintln(out, "No applications found.")
		return
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("ID", "NAME", "VERSIONS", "LATEST VERSION")

	for _, sw := range softwareList {
		latest := "None"
		if len(sw.Versions) > 0 {
			latest = sw.Versions[len(sw.Versions) - 1].Number
		}

		table.AddRow(
			sw.Identifier,
			sw.Name,
			len(sw.Versions),
			latest,
		)
	}

	fmt.Fprintln(out, table)
}