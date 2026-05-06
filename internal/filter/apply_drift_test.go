package filter_test

import (
	"testing"

	"github.com/user/driftctl-diff/internal/diff"
	"github.com/user/driftctl-diff/internal/filter"
)

func makeDriftResult() diff.Result {
	return diff.Result{
		Added: map[string]interface{}{
			"aws_instance.web":    map[string]interface{}{"ami": "ami-123"},
			"aws_s3_bucket.logs": map[string]interface{}{"bucket": "my-logs"},
		},
		Removed: map[string]interface{}{
			"aws_instance.old": map[string]interface{}{"ami": "ami-old"},
		},
		Changed: map[string]diff.Changes{
			"aws_instance.db": {{Attribute: "instance_type", From: "t2.micro", To: "t3.small"}},
		},
	}
}

func TestDriftResult_NoFilter(t *testing.T) {
	r := makeDriftResult()
	out := filter.DriftResult(r, filter.Options{})
	if len(out.Added) != 2 || len(out.Removed) != 1 || len(out.Changed) != 1 {
		t.Fatalf("expected all resources to pass, got added=%d removed=%d changed=%d",
			len(out.Added), len(out.Removed), len(out.Changed))
	}
}

func TestDriftResult_IncludeType(t *testing.T) {
	r := makeDriftResult()
	opts := filter.Options{IncludeTypes: []string{"aws_instance"}}
	out := filter.DriftResult(r, opts)

	if len(out.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(out.Added))
	}
	if _, ok := out.Added["aws_instance.web"]; !ok {
		t.Fatal("expected aws_instance.web in added")
	}
	if len(out.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(out.Changed))
	}
}

func TestDriftResult_ExcludeType(t *testing.T) {
	r := makeDriftResult()
	opts := filter.Options{ExcludeTypes: []string{"aws_s3_bucket"}}
	out := filter.DriftResult(r, opts)

	if _, ok := out.Added["aws_s3_bucket.logs"]; ok {
		t.Fatal("aws_s3_bucket.logs should have been excluded")
	}
	if len(out.Added) != 1 {
		t.Fatalf("expected 1 added after exclusion, got %d", len(out.Added))
	}
}
