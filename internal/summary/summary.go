package summary

import (
	"fmt"
	"sort"

	"github.com/your-org/driftctl-diff/internal/diff"
)

// TypeSummary holds drift counts for a single resource type.
type TypeSummary struct {
	Type    string
	Added   int
	Removed int
	Changed int
	Total   int
}

// Report holds aggregated drift statistics across all resource types.
type Report struct {
	ByType      []TypeSummary
	TotalAdded  int
	TotalRemoved int
	TotalChanged int
	TotalDrift  int
	Clean       bool
}

// Build produces a Report from a diff.Result.
func Build(result diff.Result) Report {
	counts := map[string]*TypeSummary{}

	accumulate := func(key string, added, removed, changed int) {
		if _, ok := counts[key]; !ok {
			counts[key] = &TypeSummary{Type: key}
		}
		s := counts[key]
		s.Added += added
		s.Removed += removed
		s.Changed += changed
		s.Total += added + removed + changed
	}

	for _, r := range result.Added {
		accumulate(r.Type, 1, 0, 0)
	}
	for _, r := range result.Removed {
		accumulate(r.Type, 0, 1, 0)
	}
	for _, r := range result.Changed {
		accumulate(r.Type, 0, 0, 1)
	}

	byType := make([]TypeSummary, 0, len(counts))
	for _, v := range counts {
		byType = append(byType, *v)
	}
	sort.Slice(byType, func(i, j int) bool {
		return byType[i].Type < byType[j].Type
	})

	totalAdded := len(result.Added)
	totalRemoved := len(result.Removed)
	totalChanged := len(result.Changed)
	totalDrift := totalAdded + totalRemoved + totalChanged

	return Report{
		ByType:       byType,
		TotalAdded:   totalAdded,
		TotalRemoved: totalRemoved,
		TotalChanged: totalChanged,
		TotalDrift:   totalDrift,
		Clean:        totalDrift == 0,
	}
}

// FormatLine returns a human-readable one-line summary for a TypeSummary.
func FormatLine(ts TypeSummary) string {
	return fmt.Sprintf("%-40s added=%-4d removed=%-4d changed=%-4d",
		ts.Type, ts.Added, ts.Removed, ts.Changed)
}
