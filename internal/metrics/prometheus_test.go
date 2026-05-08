package metrics_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/metrics"
)

func TestWritePrometheus_NoDrift(t *testing.T) {
	start := time.Now()
	end := start.Add(5 * time.Millisecond)
	m := metrics.Collect(diff.Result{}, start, end, 0)

	var buf bytes.Buffer
	if err := metrics.WritePrometheus(&buf, m, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, name := range []string{
		"driftctl_diff_resources_total",
		"driftctl_diff_added_total",
		"driftctl_diff_removed_total",
		"driftctl_diff_changed_total",
		"driftctl_diff_duration_seconds",
	} {
		if !strings.Contains(out, name) {
			t.Errorf("output missing metric %q", name)
		}
	}
}

func TestWritePrometheus_WithLabels(t *testing.T) {
	start := time.Now()
	end := start.Add(10 * time.Millisecond)
	m := metrics.Collect(diff.Result{}, start, end, 0)

	labels := []metrics.Label{
		{Key: "env", Value: "prod"},
		{Key: "app", Value: "myapp"},
	}
	var buf bytes.Buffer
	if err := metrics.WritePrometheus(&buf, m, labels); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `app="myapp"`) {
		t.Errorf("expected label app in output: %s", out)
	}
	if !strings.Contains(out, `env="prod"`) {
		t.Errorf("expected label env in output: %s", out)
	}
}

func TestWritePrometheus_Values(t *testing.T) {
	start := time.Now()
	end := start.Add(100 * time.Millisecond)
	changed := map[string]diff.AttributeChanges{"res.a": {}}
	m := metrics.Collect(makeResult(2, 1, changed), start, end, 4)

	var buf bytes.Buffer
	_ = metrics.WritePrometheus(&buf, m, nil)
	out := buf.String()

	for _, check := range []string{
		"driftctl_diff_added_total 2",
		"driftctl_diff_removed_total 1",
		"driftctl_diff_changed_total 1",
		"driftctl_diff_filtered_total 4",
		"driftctl_diff_resources_total 4",
	} {
		if !strings.Contains(out, check) {
			t.Errorf("expected %q in output:\n%s", check, out)
		}
	}
	_ = fmt.Sprintf("") // suppress unused import
}
