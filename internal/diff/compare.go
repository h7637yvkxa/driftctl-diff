package diff

import (
	"fmt"
	"sort"

	"github.com/user/driftctl-diff/internal/state"
)

// DriftType represents the kind of drift detected.
type DriftType string

const (
	DriftAdded   DriftType = "ADDED"
	DriftRemoved DriftType = "REMOVED"
	DriftChanged DriftType = "CHANGED"
)

// DriftEntry describes a single drift item between two states.
type DriftEntry struct {
	Key      string
	Drift    DriftType
	BaseVal  string
	TargetVal string
}

// Result holds the full diff between two state files.
type Result struct {
	Entries []DriftEntry
}

// HasDrift returns true when at least one drift entry exists.
func (r *Result) HasDrift() bool {
	return len(r.Entries) > 0
}

// Compare diffs two indexed resource maps and returns a Result.
func Compare(base, target map[string]state.Resource) *Result {
	result := &Result{}

	for key, baseRes := range base {
		if targetRes, ok := target[key]; !ok {
			result.Entries = append(result.Entries, DriftEntry{
				Key:     key,
				Drift:   DriftRemoved,
				BaseVal: fmt.Sprintf("%v", baseRes.Values),
			})
		} else if fmt.Sprintf("%v", baseRes.Values) != fmt.Sprintf("%v", targetRes.Values) {
			result.Entries = append(result.Entries, DriftEntry{
				Key:       key,
				Drift:     DriftChanged,
				BaseVal:   fmt.Sprintf("%v", baseRes.Values),
				TargetVal: fmt.Sprintf("%v", targetRes.Values),
			})
		}
	}

	for key, targetRes := range target {
		if _, ok := base[key]; !ok {
			result.Entries = append(result.Entries, DriftEntry{
				Key:       key,
				Drift:     DriftAdded,
				TargetVal: fmt.Sprintf("%v", targetRes.Values),
			})
		}
	}

	sort.Slice(result.Entries, func(i, j int) bool {
		return result.Entries[i].Key < result.Entries[j].Key
	})

	return result
}
