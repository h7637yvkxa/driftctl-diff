package lint

import (
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func makeResult(added, removed map[string]map[string]interface{}, changed map[string]map[string]diff.AttributeDiff) diff.Result {
	return diff.Result{
		Added:   added,
		Removed: removed,
		Changed: changed,
	}
}

func TestRun_NoDrift(t *testing.T) {
	result := makeResult(nil, nil, nil)
	findings := Run(result, DefaultRules())
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestRun_Added(t *testing.T) {
	result := makeResult(
		map[string]map[string]interface{}{"aws_instance.web": {"id": "i-123"}},
		nil, nil,
	)
	findings := Run(result, DefaultRules())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "DRIFT001" {
		t.Errorf("expected DRIFT001, got %s", findings[0].RuleID)
	}
	if findings[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestRun_Removed(t *testing.T) {
	result := makeResult(
		nil,
		map[string]map[string]interface{}{"aws_s3_bucket.data": {"bucket": "my-bucket"}},
		nil,
	)
	findings := Run(result, DefaultRules())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "DRIFT002" {
		t.Errorf("expected DRIFT002, got %s", findings[0].RuleID)
	}
	if findings[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestRun_Changed(t *testing.T) {
	result := makeResult(nil, nil,
		map[string]map[string]diff.AttributeDiff{
			"aws_instance.web": {
				"instance_type": {Old: "t2.micro", New: "t3.small"},
			},
		},
	)
	findings := Run(result, DefaultRules())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "DRIFT003" {
		t.Errorf("expected DRIFT003, got %s", findings[0].RuleID)
	}
}

func TestFinding_String(t *testing.T) {
	f := Finding{RuleID: "DRIFT001", Severity: "warning", Resource: "aws_instance.web", Detail: "added"}
	s := f.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
