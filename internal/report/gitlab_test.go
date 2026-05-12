package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func decodeGitLab(t *testing.T, buf *bytes.Buffer) []gitLabIssue {
	t.Helper()
	var issues []gitLabIssue
	if err := json.NewDecoder(buf).Decode(&issues); err != nil {
		t.Fatalf("decode gitlab: %v", err)
	}
	return issues
}

func TestWriteGitLab_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := diff.Result{}
	if err := writeGitLab(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	issues := decodeGitLab(t, &buf)
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestWriteGitLab_Added(t *testing.T) {
	var buf bytes.Buffer
	r := diff.Result{
		Added: map[string]diff.Resource{
			"aws_instance.web": {Type: "aws_instance", Name: "web"},
		},
	}
	if err := writeGitLab(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	issues := decodeGitLab(t, &buf)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Description, "aws_instance.web") {
		t.Errorf("expected description to mention resource key, got %q", issues[0].Description)
	}
	if issues[0].Severity != "minor" {
		t.Errorf("expected severity minor, got %q", issues[0].Severity)
	}
	if !strings.HasPrefix(issues[0].Fingerprint, "added-") {
		t.Errorf("expected fingerprint prefix 'added-', got %q", issues[0].Fingerprint)
	}
}

func TestWriteGitLab_Removed(t *testing.T) {
	var buf bytes.Buffer
	r := diff.Result{
		Removed: map[string]diff.Resource{
			"aws_s3_bucket.data": {Type: "aws_s3_bucket", Name: "data"},
		},
	}
	if err := writeGitLab(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	issues := decodeGitLab(t, &buf)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "major" {
		t.Errorf("expected severity major for removed resource, got %q", issues[0].Severity)
	}
}

func TestWriteGitLab_Changed(t *testing.T) {
	var buf bytes.Buffer
	r := diff.Result{
		Changed: map[string]diff.Change{
			"aws_instance.app": {
				Resource: diff.Resource{Type: "aws_instance", Name: "app"},
				Attributes: map[string]diff.AttrDiff{
					"instance_type": {Old: "t2.micro", New: "t3.small"},
					"ami":           {Old: "ami-aaa", New: "ami-bbb"},
				},
			},
		},
	}
	if err := writeGitLab(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	issues := decodeGitLab(t, &buf)
	if len(issues) != 2 {
		t.Errorf("expected 2 issues (one per attribute), got %d", len(issues))
	}
	for _, iss := range issues {
		if !strings.HasPrefix(iss.Fingerprint, "changed-") {
			t.Errorf("expected fingerprint prefix 'changed-', got %q", iss.Fingerprint)
		}
	}
}
