package diff

// Change describes attribute-level differences for a single resource.
type Change struct {
	// Key is the resource key (type.name).
	Key string
	// Attributes maps attribute name to [before, after] values.
	Attributes map[string][2]string
}

// Result holds the full drift comparison output.
type Result struct {
	// Added contains resource keys present in the target but not the source.
	Added []string
	// Removed contains resource keys present in the source but not the target.
	Removed []string
	// Changed contains resources present in both states with differing attributes.
	Changed []Change
}

// HasDrift returns true when any drift was detected.
func (r *Result) HasDrift() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// Summary returns a human-readable one-line summary of the result.
func (r *Result) Summary() string {
	if !r.HasDrift() {
		return "no drift detected"
	}
	return formatSummary(len(r.Added), len(r.Removed), len(r.Changed))
}

func formatSummary(added, removed, changed int) string {
	return (
		"drift detected: " +
			plural(added, "added") + ", " +
			plural(removed, "removed") + ", " +
			plural(changed, "changed")
	)
}

func plural(n int, label string) string {
	if n == 1 {
		return "1 " + label
	}
	return fmt.Sprintf("%d %s", n, label)
}
