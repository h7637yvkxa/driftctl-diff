package diff

import (
	"testing"

	"github.com/user/driftctl-diff/internal/state"
)

func makeIndex(resources map[string]map[string]interface{}) map[string]state.Resource {
	idx := make(map[string]state.Resource, len(resources))
	for k, vals := range resources {
		idx[k] = state.Resource{Values: vals}
	}
	return idx
}

func TestCompare_NoDrift(t *testing.T) {
	base := makeIndex(map[string]map[string]interface{}{
		"aws_instance.web": {"ami": "ami-123", "instance_type": "t2.micro"},
	})
	target := makeIndex(map[string]map[string]interface{}{
		"aws_instance.web": {"ami": "ami-123", "instance_type": "t2.micro"},
	})

	result := Compare(base, target)
	if result.HasDrift() {
		t.Errorf("expected no drift, got %d entries", len(result.Entries))
	}
}

func TestCompare_Added(t *testing.T) {
	base := makeIndex(map[string]map[string]interface{}{})
	target := makeIndex(map[string]map[string]interface{}{
		"aws_s3_bucket.logs": {"bucket": "my-logs"},
	})

	result := Compare(base, target)
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Drift != DriftAdded {
		t.Errorf("expected ADDED, got %s", result.Entries[0].Drift)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := makeIndex(map[string]map[string]interface{}{
		"aws_s3_bucket.logs": {"bucket": "my-logs"},
	})
	target := makeIndex(map[string]map[string]interface{}{})

	result := Compare(base, target)
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Drift != DriftRemoved {
		t.Errorf("expected REMOVED, got %s", result.Entries[0].Drift)
	}
}

func TestCompare_Changed(t *testing.T) {
	base := makeIndex(map[string]map[string]interface{}{
		"aws_instance.web": {"instance_type": "t2.micro"},
	})
	target := makeIndex(map[string]map[string]interface{}{
		"aws_instance.web": {"instance_type": "t3.medium"},
	})

	result := Compare(base, target)
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Drift != DriftChanged {
		t.Errorf("expected CHANGED, got %s", result.Entries[0].Drift)
	}
}
