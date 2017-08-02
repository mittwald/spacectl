package helper

import (
	"time"
	"math"
	"fmt"
	"strings"
)

func HumanReadableDateDiff(t1, t2 time.Time) string {
	diff := t1.Round(time.Second).Sub(t2.Round(time.Second))
	components := []string{}

	if diff >= (7 * 24 * time.Hour) {
		weeks := math.Floor(float64(diff / (7 * 24 * time.Hour)))
		diff = diff % (7 * 24 * time.Hour)
		components = append(components, fmt.Sprintf("%dw", int(weeks)))
	}

	if diff >= (24 * time.Hour) {
		days := math.Floor(float64(diff / (24 * time.Hour)))
		diff = diff % (24 * time.Hour)
		components = append(components, fmt.Sprintf("%dd", int(days)))
	}

	if diff >= (1 * time.Hour) {
		hours := math.Floor(float64(diff / (1 * time.Hour)))
		diff = diff % (1 * time.Hour)
		components = append(components, fmt.Sprintf("%dh", int(hours)))
	}

	if diff >= (1 * time.Minute) {
		minutes := math.Floor(float64(diff / (1 * time.Minute)))
		diff = diff % (1 * time.Minute)
		components = append(components, fmt.Sprintf("%dm", int(minutes)))
	}

	components = append(components, fmt.Sprintf("%ds", diff / time.Second))

	return strings.Join(components[0:2], "")
}