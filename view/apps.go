package view

import (
	"fmt"
	"io"

	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/software"
)

type AppsView interface {
	List(softwareList []software.Software, out io.Writer)
}

type TabularAppsView struct{}

func (t TabularAppsView) List(softwareList []software.Software, out io.Writer) {
	if len(softwareList) == 0 {
		fmt.Fprintln(out, "No applications found.")
		return
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("ID", "NAME", "VERSIONS", "LATEST VERSION", "TYPE")

	for _, sw := range softwareList {
		latest := "None"
		if len(sw.Versions) > 0 {
			latest = sw.Versions[len(sw.Versions)-1].Number
		}

		table.AddRow(
			sw.Identifier,
			sw.Name,
			len(sw.Versions),
			latest,
			sw.Type,
		)
	}

	fmt.Fprintln(out, table)
}

type AppVersionView interface {
	List(appID string, appName string, versionList []software.Version, out io.Writer)
}

type TabularAppVersionView struct{}

func (t TabularAppVersionView) List(appID string, appName string, versionList []software.Version, out io.Writer) {
	if appID == "" || appName == "" {
		fmt.Fprintln(out, "No application found.")
		return
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("ID", "NAME", "VERSIONS")

	table.AddRow(
		appID,
		appName,
		versionList[0].Number,
	)
	versionList = versionList[1:]

	for _, sw := range versionList {
		table.AddRow("", "", sw.Number)
	}

	fmt.Fprintln(out, table)
}
