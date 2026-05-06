package baseline_test

import (
	"testing"

	"github.com/owner/driftctl-diff/internal/baseline"
	"github.com/owner/driftctl-diff/internal/diff"
)

func makeDiff() diff.DiffResult {
	return diff.DiffResult{
		Added: []diff.ResourceDiff{
			{Type: "aws_instance", Name: "web"},
		},
		Removed: []diff.ResourceDiff{
			{Type: "aws_security_group", Name: "default"},
		},
		Changed: []diff.ResourceDiff{
			{
				Type: "aws_s3_bucket",
				Name: "assets",
				Attributes: []diff.AttributeDiff{
					{Key: "tags", OldValue: "a", NewValue: "b"},
					{Key: "region", OldValue: "us-east-1", NewValue: "eu-west-1"},
				},
			},
		},
	}
}

func TestApply_NilBaseline(t *testing.T) {
	result := makeDiff()
	out, n := baseline.Apply(result, nil)
	if n != 0 {
		t.Errorf("expected 0 suppressed, got %d", n)
	}
	if len(out.Added) != 1 || len(out.Removed) != 1 || len(out.Changed) != 1 {
		t.Error("result should be unchanged")
	}
}

func TestApply_SuppressAddedResource(t *testing.T) {
	b := &baseline.Baseline{
		Entries: []baseline.BaselineEntry{
			{ResourceType: "aws_instance", ResourceName: "web"},
		},
	}
	out, n := baseline.Apply(makeDiff(), b)
	if n != 1 {
		t.Errorf("expected 1 suppressed, got %d", n)
	}
	if len(out.Added) != 0 {
		t.Errorf("expected Added to be empty, got %d", len(out.Added))
	}
}

func TestApply_SuppressAttribute(t *testing.T) {
	b := &baseline.Baseline{
		Entries: []baseline.BaselineEntry{
			{ResourceType: "aws_s3_bucket", ResourceName: "assets", Attribute: "tags"},
		},
	}
	out, n := baseline.Apply(makeDiff(), b)
	if n != 1 {
		t.Errorf("expected 1 suppressed, got %d", n)
	}
	if len(out.Changed) != 1 {
		t.Fatalf("expected 1 changed resource, got %d", len(out.Changed))
	}
	if len(out.Changed[0].Attributes) != 1 {
		t.Errorf("expected 1 remaining attribute, got %d", len(out.Changed[0].Attributes))
	}
}

func TestApply_SuppressWholeChangedResource(t *testing.T) {
	b := &baseline.Baseline{
		Entries: []baseline.BaselineEntry{
			{ResourceType: "aws_s3_bucket", ResourceName: "assets"},
		},
	}
	out, n := baseline.Apply(makeDiff(), b)
	if n != 1 {
		t.Errorf("expected 1 suppressed, got %d", n)
	}
	if len(out.Changed) != 0 {
		t.Errorf("expected Changed to be empty, got %d", len(out.Changed))
	}
}
