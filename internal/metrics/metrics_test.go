package metrics_test

import (
	"strings"
	"testing"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/metrics"
)

func makeResult(added, removed int, changed map[string]diff.AttributeChanges) diff.Result {
	r := diff.Result{
		Added:   make(map[string]diff.Resource),
		Removed: make(map[string]diff.Resource),
		Changed: changed,
	}
	for i := 0; i < added; i++ {
		k := fmt.Sprintf("aws_instance.a%d", i)
		r.Added[k] = diff.Resource{}
	}
	for i := 0; i < removed; i++ {
		k := fmt.Sprintf("aws_instance.r%d", i)
		r.Removed[k] = diff.Resource{}
	}
	return r
}

func TestCollect_NoDrift(t *testing.T) {
	start := time.Now()
	end := start.Add(10 * time.Millisecond)
	m := metrics.Collect(makeResult(0, 0, nil), start, end, 0)
	if m.HasDrift() {
		t.Fatal("expected no drift")
	}
	if m.TotalResources != 0 {
		t.Fatalf("expected 0 total, got %d", m.TotalResources)
	}
}

func TestCollect_WithDrift(t *testing.T) {
	start := time.Now()
	end := start.Add(50 * time.Millisecond)
	changed := map[string]diff.AttributeChanges{"aws_s3_bucket.b": {}}
	m := metrics.Collect(makeResult(2, 1, changed), start, end, 3)
	if !m.HasDrift() {
		t.Fatal("expected drift")
	}
	if m.Added != 2 || m.Removed != 1 || m.Changed != 1 {
		t.Fatalf("unexpected counts: %+v", m)
	}
	if m.Filtered != 3 {
		t.Fatalf("expected filtered=3, got %d", m.Filtered)
	}
	if m.TotalResources != 4 {
		t.Fatalf("expected total=4, got %d", m.TotalResources)
	}
}

func TestSummary_ContainsFields(t *testing.T) {
	start := time.Now()
	end := start.Add(25 * time.Millisecond)
	m := metrics.Collect(makeResult(1, 0, nil), start, end, 2)
	s := m.Summary()
	for _, want := range []string{"added", "removed", "changed", "filtered", "duration"} {
		if !strings.Contains(s, want) {
			t.Errorf("summary missing %q: %s", want, s)
		}
	}
}
