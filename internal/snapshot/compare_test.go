package snapshot_test

import (
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/snapshot"
)

func makeSnap(result diff.Result) *snapshot.Snapshot {
	return &snapshot.Snapshot{Result: result}
}

func TestCompareTo_AllNew(t *testing.T) {
	snap := makeSnap(diff.Result{})
	current := makeResult()
	delta := snap.CompareTo(current)

	if len(delta.New) != 3 {
		t.Errorf("New = %d, want 3", len(delta.New))
	}
	if len(delta.Resolved) != 0 {
		t.Errorf("Resolved = %d, want 0", len(delta.Resolved))
	}
	if len(delta.Persistent) != 0 {
		t.Errorf("Persistent = %d, want 0", len(delta.Persistent))
	}
}

func TestCompareTo_AllResolved(t *testing.T) {
	snap := makeSnap(makeResult())
	delta := snap.CompareTo(diff.Result{})

	if len(delta.Resolved) != 3 {
		t.Errorf("Resolved = %d, want 3", len(delta.Resolved))
	}
	if len(delta.New) != 0 || len(delta.Persistent) != 0 {
		t.Errorf("unexpected New/Persistent: %v / %v", delta.New, delta.Persistent)
	}
}

func TestCompareTo_Persistent(t *testing.T) {
	prev := diff.Result{Added: []diff.Resource{{Key: "aws_s3_bucket.logs"}}}
	snap := makeSnap(prev)
	current := diff.Result{Added: []diff.Resource{{Key: "aws_s3_bucket.logs"}}}
	delta := snap.CompareTo(current)

	if len(delta.Persistent) != 1 {
		t.Errorf("Persistent = %d, want 1", len(delta.Persistent))
	}
	if len(delta.New) != 0 || len(delta.Resolved) != 0 {
		t.Errorf("unexpected New/Resolved: %v / %v", delta.New, delta.Resolved)
	}
}

func TestDelta_Summary(t *testing.T) {
	d := snapshot.Delta{New: []string{"a"}, Resolved: []string{"b", "c"}, Persistent: []string{}}
	s := d.Summary()
	if !strings.Contains(s, "1 new") || !strings.Contains(s, "2 resolved") {
		t.Errorf("unexpected summary: %q", s)
	}
}
