package snapshot

import (
	"fmt"

	"github.com/owner/driftctl-diff/internal/diff"
)

// Delta describes how a diff result has changed relative to a previous snapshot.
type Delta struct {
	New        []string // resource keys newly drifted
	Resolved   []string // resource keys no longer drifting
	Persistent []string // resource keys still drifting
}

// CompareTo computes the delta between a saved snapshot and a current diff result.
func (s *Snapshot) CompareTo(current diff.Result) Delta {
	prev := keySet(s.Result)
	curr := keySet(current)

	var d Delta
	for k := range curr {
		if prev[k] {
			d.Persistent = append(d.Persistent, k)
		} else {
			d.New = append(d.New, k)
		}
	}
	for k := range prev {
		if !curr[k] {
			d.Resolved = append(d.Resolved, k)
		}
	}
	return d
}

// Summary returns a human-readable summary of the delta.
func (d Delta) Summary() string {
	return fmt.Sprintf(
		"snapshot delta: %d new drift(s), %d resolved, %d persistent",
		len(d.New), len(d.Resolved), len(d.Persistent),
	)
}

func keySet(r diff.Result) map[string]bool {
	keys := make(map[string]bool)
	for _, res := range r.Added {
		keys[res.Key] = true
	}
	for _, res := range r.Removed {
		keys[res.Key] = true
	}
	for _, ch := range r.Changed {
		keys[ch.Key] = true
	}
	return keys
}
