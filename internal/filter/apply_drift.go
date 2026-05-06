package filter

import (
	"github.com/user/driftctl-diff/internal/diff"
)

// DriftResult filters a diff.Result by the given Options, returning a new
// Result containing only the resources that pass the filter.
func DriftResult(result diff.Result, opts Options) diff.Result {
	filtered := diff.Result{
		Added:   make(map[string]interface{}),
		Removed: make(map[string]interface{}),
		Changed: make(map[string]diff.Changes),
	}

	for k, v := range result.Added {
		if Apply(k, opts) {
			filtered.Added[k] = v
		}
	}

	for k, v := range result.Removed {
		if Apply(k, opts) {
			filtered.Removed[k] = v
		}
	}

	for k, v := range result.Changed {
		if Apply(k, opts) {
			filtered.Changed[k] = v
		}
	}

	return filtered
}
