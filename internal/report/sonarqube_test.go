package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func decodeSonar(t *testing.T, buf *bytes.Buffer) sonarReport {
	t.Helper()
	var r sonarReport
	if err := json.NewDecoder(buf).Decode(&r); err != nil {
		t.Fatalf("decode sonar: %v", err)
	}
	return r
}

func TestWriteSonarQube_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}
	if err := writeSonarQube(&buf, result, "state.tfstate"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := decodeSonar(t, &buf)
	if len(r.Issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(r.Issues))
	}
}

func TestWriteSonarQube_Added(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Added: map[string]diff.Resource{
			"aws_s3_bucket.logs": {Type: "aws_s3_bucket", Name: "logs"},
		},
	}
	if err := writeSonarQube(&buf, result, "prod.tfstate"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := decodeSonar(t, &buf)
	if len(r.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(r.Issues))
	}
	issue := r.Issues[0]
	if issue.RuleID != "drift.resource.added" {
		t.Errorf("ruleId = %q, want drift.resource.added", issue.RuleID)
	}
	if issue.Severity != "MAJOR" {
		t.Errorf("severity = %q, want MAJOR", issue.Severity)
	}
	if issue.PrimaryLocation.FilePath != "prod.tfstate" {
		t.Errorf("filePath = %q, want prod.tfstate", issue.PrimaryLocation.FilePath)
	}
}

func TestWriteSonarQube_Removed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Removed: map[string]diff.Resource{
			"aws_instance.web": {Type: "aws_instance", Name: "web"},
		},
	}
	if err := writeSonarQube(&buf, result, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := decodeSonar(t, &buf)
	if len(r.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(r.Issues))
	}
	if r.Issues[0].RuleID != "drift.resource.removed" {
		t.Errorf("ruleId = %q, want drift.resource.removed", r.Issues[0].RuleID)
	}
	if r.Issues[0].Severity != "CRITICAL" {
		t.Errorf("severity = %q, want CRITICAL", r.Issues[0].Severity)
	}
	if r.Issues[0].PrimaryLocation.FilePath != "terraform.tfstate" {
		t.Errorf("default filePath not applied")
	}
}

func TestWriteSonarQube_Changed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Changed: map[string]map[string]diff.AttributeDiff{
			"aws_security_group.default": {
				"ingress": {Old: "0.0.0.0/0", New: "10.0.0.0/8"},
			},
		},
	}
	if err := writeSonarQube(&buf, result, "dev.tfstate"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := decodeSonar(t, &buf)
	if len(r.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(r.Issues))
	}
	if r.Issues[0].RuleID != "drift.resource.changed" {
		t.Errorf("ruleId = %q, want drift.resource.changed", r.Issues[0].RuleID)
	}
	if r.Issues[0].Type != "BUG" {
		t.Errorf("type = %q, want BUG", r.Issues[0].Type)
	}
	if r.Issues[0].EngineID != "driftctl-diff" {
		t.Errorf("engineId = %q, want driftctl-diff", r.Issues[0].EngineID)
	}
}
