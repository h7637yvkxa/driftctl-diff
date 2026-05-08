package metrics

import (
	"fmt"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
)

// RunMetrics holds timing and count statistics for a single diff run.
type RunMetrics struct {
	StartedAt     time.Time
	FinishedAt    time.Time
	Duration      time.Duration
	TotalResources int
	Added         int
	Removed       int
	Changed       int
	Filtered      int
}

// Collect builds a RunMetrics value from a completed DiffResult.
func Collect(result diff.Result, start, end time.Time, filtered int) RunMetrics {
	m := RunMetrics{
		StartedAt:  start,
		FinishedAt: end,
		Duration:   end.Sub(start),
		Added:      len(result.Added),
		Removed:    len(result.Removed),
		Changed:    len(result.Changed),
		Filtered:   filtered,
	}
	m.TotalResources = m.Added + m.Removed + m.Changed
	return m
}

// HasDrift returns true when any drift was detected.
func (m RunMetrics) HasDrift() bool {
	return m.TotalResources > 0
}

// Summary returns a single-line human-readable summary.
func (m RunMetrics) Summary() string {
	return fmt.Sprintf(
		"drift: %d added, %d removed, %d changed | filtered: %d | duration: %s",
		m.Added, m.Removed, m.Changed, m.Filtered, m.Duration.Round(time.Millisecond),
	)
}
