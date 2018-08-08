package view

import (
	"time"

	"fmt"

	ui "github.com/gizak/termui"
	"github.com/mittwald/spacectl/client/spaces"
)

type SpaceMetricsView struct {
}

func (*SpaceMetricsView) Render(scope string, loadMetrics func() (spaces.ComputeMetricPointList, error)) error {

	err := ui.Init()
	if err != nil {
		return err
	}

	metrics, err := loadMetrics()
	if err != nil {
		return err
	}

	// "Mon, 02 Jan 2006 15:04:05 MST"
	format := time.Kitchen
	if scope == spaces.ScopeWeek {
		format = "02 Jan 15:04"
	} else if scope == spaces.ScopeMonth {
		format = "02 Jan"
	} else if scope == spaces.ScopeYear {
		format = "02 Jan"
	}

	cpuLast := metrics[len(metrics)-1].CPU
	memLast := metrics[len(metrics)-1].Memory

	cpuChart := ui.NewLineChart()
	cpuChart.BorderLabel = "CPU usage (CPU seconds per second)"
	cpuChart.Data = metrics.CPUUsage()
	cpuChart.DataLabels = metrics.DateStrings(format)
	cpuChart.Height = 12

	cpuTable := ui.NewTable()
	cpuTable.BorderLabel = "CPU usage (sum of all components)"
	cpuTable.Height = 5
	cpuTable.Separator = false
	cpuTable.Rows = [][]string{
		{"Usage", fmt.Sprintf("%.2f", cpuLast.Usage)},
		{"Guaranteed", fmt.Sprintf("%.2f", cpuLast.Request)},
		{"Burst", fmt.Sprintf("%.2f", cpuLast.Limit)},
	}

	memChart := ui.NewLineChart()
	memChart.BorderLabel = "Memory usage (MiB)"
	memChart.Data = metrics.MemoryUsage(spaces.UnitMB)
	memChart.DataLabels = metrics.DateStrings(format)
	memChart.Height = 12

	memTable := ui.NewTable()
	memTable.BorderLabel = "Memory usage (sum of all components)"
	memTable.Height = 5
	memTable.Separator = false
	memTable.Rows = [][]string{
		{"Usage", fmt.Sprintf("%.2f MiB", memLast.Usage/spaces.UnitMB)},
		{"Guaranteed", fmt.Sprintf("%.2f MiB", memLast.Request/spaces.UnitMB)},
		{"Burst", fmt.Sprintf("%.2f MiB", memLast.Limit/spaces.UnitMB)},
	}

	lastUpdated := ui.NewPar("Last updated on " + time.Now().Format("15:04:05 MST"))
	lastUpdated.Border = false
	lastUpdated.Height = 2

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, &ui.Par{Block: ui.Block{Border: false, Height: 1}, Text: "SPACES metrics"}),
		),
		ui.NewRow(
			ui.NewCol(12, 0, lastUpdated),
		),
		ui.NewRow(
			ui.NewCol(6, 0, cpuChart),
			ui.NewCol(6, 0, memChart),
		),
		ui.NewRow(
			ui.NewCol(6, 0, cpuTable),
			ui.NewCol(6, 0, memTable),
		),
	)
	ui.Body.Align()

	ui.Handle("/sys/kbd/q", func(event ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/r", func(event ui.Event) {
		metrics, _ = loadMetrics()
		lastUpdated.Text = "Last updated on " + time.Now().Format("15:04:05 MST")

		cpuChart.Data = metrics.CPUUsage()
		cpuChart.DataLabels = metrics.DateStrings(format)
		cpuTable.Rows = [][]string{
			{"Usage", fmt.Sprintf("%.2f", cpuLast.Usage)},
			{"Guaranteed", fmt.Sprintf("%.2f", cpuLast.Request)},
			{"Burst", fmt.Sprintf("%.2f", cpuLast.Limit)},
		}

		memChart.Data = metrics.MemoryUsage(spaces.UnitMB)
		memChart.DataLabels = metrics.DateStrings(format)
		memTable.Rows = [][]string{
			{"Usage", fmt.Sprintf("%.2f MiB", memLast.Usage/spaces.UnitMB)},
			{"Guaranteed", fmt.Sprintf("%.2f MiB", memLast.Request/spaces.UnitMB)},
			{"Burst", fmt.Sprintf("%.2f MiB", memLast.Limit/spaces.UnitMB)},
		}

		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kbd/C-c", func(event ui.Event) {
		ui.StopLoop()
	})

	ui.Render(ui.Body)
	ui.Loop()
	ui.Close()

	return nil
}
