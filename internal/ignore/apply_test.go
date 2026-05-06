package ignore

import (
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func makeResult() diff.Result {
	return diff.Result{
		Added: []diff.Resource{
			{Type: "aws_s3_bucket", Name: "logs"},
		},
		Removed: []diff.Resource{
			{Type: "aws_instance", Name: "web"},
		},
		Changed: []diff.Change{
			{
				Type: "aws_instance",
				Name: "app",
				Attributes: []diff.AttributeDiff{
					{Key: "tags", OldValue: "a", NewValue: "b"},
					{Key: "ami", OldValue: "old", NewValue: "new"},
				},
			},
		},
	}
}

func TestFilterResult_NilSet(t *testing.T) {
	r := makeResult()
	out := FilterResult(r, nil)
	if len(out.Added) != 1 || len(out.Removed) != 1 || len(out.Changed) != 1 {
		t.Error("nil set should not filter anything")
	}
}

func TestFilterResult_IgnoreAddedResource(t *testing.T) {
	set := &Set{rules: []Rule{{ResourceType: "aws_s3_bucket", ResourceName: "logs"}}}
	out := FilterResult(makeResult(), set)
	if len(out.Added) != 0 {
		t.Errorf("expected added resource to be ignored, got %d", len(out.Added))
	}
}

func TestFilterResult_IgnoreAttribute(t *testing.T) {
	set := &Set{rules: []Rule{{ResourceType: "aws_instance", ResourceName: "app", Attribute: "tags"}}}
	out := FilterResult(makeResult(), set)
	if len(out.Changed) != 1 {
		t.Fatal("expected one changed resource")
	}
	if len(out.Changed[0].Attributes) != 1 || out.Changed[0].Attributes[0].Key != "ami" {
		t.Error("expected only 'ami' attribute to remain")
	}
}

func TestFilterResult_IgnoreWholeChangedResource(t *testing.T) {
	set := &Set{rules: []Rule{{ResourceType: "aws_instance", ResourceName: "app"}}}
	out := FilterResult(makeResult(), set)
	if len(out.Changed) != 0 {
		t.Errorf("expected changed resource to be fully ignored, got %d", len(out.Changed))
	}
}
