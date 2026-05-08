package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func TestWriteSARIF_NoDrift(t *testing.T) {
	result := diff.Result{}
	var buf bytes.Buffer
	if err := writeSARIF(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var log sarifLog
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if log.Version != "2.1.0" {
		t.Errorf("expected version 2.1.0, got %q", log.Version)
	}
	if len(log.Runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(log.Runs))
	}
	if len(log.Runs[0].Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(log.Runs[0].Results))
	}
}

func TestWriteSARIF_Added(t *testing.T) {
	result := diff.Result{
		Added: map[string]diff.Resource{
			"aws_s3_bucket.logs": {Type: "aws_s3_bucket", Name: "logs"},
		},
	}
	var buf bytes.Buffer
	if err := writeSARIF(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var log sarifLog
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := log.Runs[0].Results
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].RuleID != "drift/added" {
		t.Errorf("expected ruleId drift/added, got %q", results[0].RuleID)
	}
	if results[0].Level != "warning" {
		t.Errorf("expected level warning, got %q", results[0].Level)
	}
	if !strings.Contains(results[0].Message.Text, "aws_s3_bucket.logs") {
		t.Errorf("message should mention resource key, got %q", results[0].Message.Text)
	}
}

func TestWriteSARIF_Changed(t *testing.T) {
	result := diff.Result{
		Changed: map[string]diff.ChangedResource{
			"aws_instance.web": {
				Type: "aws_instance",
				Name: "web",
				Attributes: map[string]diff.AttributeDiff{
					"instance_type": {Source: "t2.micro", Target: "t3.small"},
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := writeSARIF(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var log sarifLog
	if err := json.Unmarshal(buf.Bytes(), &log); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	results := log.Runs[0].Results
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].RuleID != "drift/changed" {
		t.Errorf("expected ruleId drift/changed, got %q", results[0].RuleID)
	}
	if results[0].Level != "note" {
		t.Errorf("expected level note, got %q", results[0].Level)
	}
	if !strings.Contains(results[0].Message.Text, "instance_type") {
		t.Errorf("message should mention attribute, got %q", results[0].Message.Text)
	}
}
