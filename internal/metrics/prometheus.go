package metrics

import (
	"fmt"
	"io"
	"sort"
)

// Label is a key=value pair for Prometheus output.
type Label struct {
	Key   string
	Value string
}

// WritePrometheus writes the RunMetrics in the Prometheus text exposition
// format to w. The optional labels slice is attached to every metric.
func WritePrometheus(w io.Writer, m RunMetrics, labels []Label) error {
	labelStr := buildLabelStr(labels)

	lines := []struct {
		help  string
		name  string
		value interface{}
	}{
		{"Total drifted resources detected", "driftctl_diff_resources_total", m.TotalResources},
		{"Resources added in target state", "driftctl_diff_added_total", m.Added},
		{"Resources removed from target state", "driftctl_diff_removed_total", m.Removed},
		{"Resources with changed attributes", "driftctl_diff_changed_total", m.Changed},
		{"Resources filtered from output", "driftctl_diff_filtered_total", m.Filtered},
		{"Duration of the diff run in seconds", "driftctl_diff_duration_seconds", m.Duration.Seconds()},
	}

	for _, l := range lines {
		if _, err := fmt.Fprintf(w, "# HELP %s %s\n# TYPE %s gauge\n%s%s %v\n",
			l.name, l.help, l.name, l.name, labelStr, l.value); err != nil {
			return err
		}
	}
	return nil
}

func buildLabelStr(labels []Label) string {
	if len(labels) == 0 {
		return ""
	}
	sort.Slice(labels, func(i, j int) bool { return labels[i].Key < labels[j].Key })
	out := "{"
	for i, l := range labels {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf(`%s=%q`, l.Key, l.Value)
	}
	return out + "}"
}
