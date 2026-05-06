package ignore

import (
	"github.com/owner/driftctl-diff/internal/diff"
)

// FilterResult removes changes that are suppressed by the ignore set.
// A resource is fully ignored if a whole-resource rule matches.
// Individual attribute changes are dropped if an attribute rule matches.
func FilterResult(result diff.Result, set *Set) diff.Result {
	if set == nil || len(set.rules) == 0 {
		return result
	}

	filtered := diff.Result{
		Added:   filterResources(result.Added, set),
		Removed: filterResources(result.Removed, set),
		Changed: filterChanged(result.Changed, set),
	}
	return filtered
}

func filterResources(resources []diff.Resource, set *Set) []diff.Resource {
	out := make([]diff.Resource, 0, len(resources))
	for _, r := range resources {
		if !set.Matches(r.Type, r.Name, "") {
			out = append(out, r)
		}
	}
	return out
}

func filterChanged(changes []diff.Change, set *Set) []diff.Change {
	out := make([]diff.Change, 0, len(changes))
	for _, c := range changes {
		if set.Matches(c.Type, c.Name, "") {
			continue // whole resource ignored
		}
		kept := make([]diff.AttributeDiff, 0, len(c.Attributes))
		for _, attr := range c.Attributes {
			if !set.Matches(c.Type, c.Name, attr.Key) {
				kept = append(kept, attr)
			}
		}
		if len(kept) > 0 {
			c.Attributes = kept
			out = append(out, c)
		}
	}
	return out
}
