package baseline

import "github.com/owner/driftctl-diff/internal/diff"

// Apply removes drift entries that are present in the baseline from the result.
// It returns a new DiffResult with baselined items suppressed and a count of
// how many entries were suppressed.
func Apply(result diff.DiffResult, b *Baseline) (diff.DiffResult, int) {
	if b == nil || len(b.Entries) == 0 {
		return result, 0
	}

	// Build a lookup set from baseline entries.
	suppressed := make(map[string]bool, len(b.Entries))
	for _, e := range b.Entries {
		suppressed[EntryKey(e)] = true
	}

	count := 0

	// Filter Added resources.
	var added []diff.ResourceDiff
	for _, r := range result.Added {
		key := r.Type + "." + r.Name
		if suppressed[key] {
			count++
			continue
		}
		added = append(added, r)
	}

	// Filter Removed resources.
	var removed []diff.ResourceDiff
	for _, r := range result.Removed {
		key := r.Type + "." + r.Name
		if suppressed[key] {
			count++
			continue
		}
		removed = append(removed, r)
	}

	// Filter Changed resources / attributes.
	var changed []diff.ResourceDiff
	for _, r := range result.Changed {
		resKey := r.Type + "." + r.Name
		if suppressed[resKey] {
			count++
			continue
		}
		var attrs []diff.AttributeDiff
		for _, a := range r.Attributes {
			attrKey := resKey + ":" + a.Key
			if suppressed[attrKey] {
				count++
				continue
			}
			attrs = append(attrs, a)
		}
		if len(attrs) > 0 {
			r.Attributes = attrs
			changed = append(changed, r)
		}
	}

	return diff.DiffResult{
		Added:   added,
		Removed: removed,
		Changed: changed,
	}, count
}
