package spaces

import "time"

const (
	UnitBytes = 1
	UnitKB    = 1024
	UnitMB    = 1024 * 1024
	UnitGB    = 1024 * 1024 * 1024
)

type MetricDateRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type MetricUsage struct {
	Usage   float64 `json:"usage"`
	Limit   float64 `json:"limit"`
	Request float64 `json:"request"`
}

type StorageUsage struct {
	Mean         float64 `json:"mean"`
	Percentile90 float64 `json:"p90"`
}

type ComputeMetricPoint struct {
	Date   MetricDateRange `json:"date"`
	Stage  string          `json:"stage"`
	CPU    MetricUsage     `json:"cpu"`
	Memory MetricUsage     `json:"memory"`
}

type ComputeMetricPointList []ComputeMetricPoint

func (l ComputeMetricPointList) Dates() []time.Time {
	results := make([]time.Time, len(l))

	for i := range l {
		results[i] = l[i].Date.From
	}

	return results
}

func (l ComputeMetricPointList) DateStrings(layout string) []string {
	results := make([]string, len(l))

	for i := range l {
		results[i] = l[i].Date.From.Format(layout)
	}

	return results
}

func (l ComputeMetricPointList) CPUUsage() []float64 {
	results := make([]float64, len(l))

	for i := range l {
		results[i] = l[i].CPU.Usage
	}

	return results
}

func (l ComputeMetricPointList) MemoryUsage(unit int) []float64 {
	results := make([]float64, len(l))

	for i := range l {
		results[i] = l[i].Memory.Usage / float64(unit)
	}

	return results
}

type StorageMetricPoint struct {
	Date     MetricDateRange `json:"date"`
	Stage    string          `json:"stage"`
	Files    StorageUsage    `json:"files"`
	Database StorageUsage    `json:"database"`
}
