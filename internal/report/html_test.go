package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftctl-diff/internal/diff"
)

func TestWriteHTML_NoDrift(t *testing.T) {
	result := diff.Result{}
	var buf bytes.Buffer
	if err := writeHTML(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "<!DOCTYPE html>") {
		t.Error("expected HTML doctype")
	}
	if strings.Contains(out, "<table>") {
		t.Error("expected no table when no drift")
	}
	if !strings.Contains(out, "no drift") && !strings.Contains(out, "No drift") {
		t.Logf("summary line: %s", out)
	}
}

func TestWriteHTML_Added(t *testing.T) {
	result := diff.Result{
		Added: []string{"aws_s3_bucket.logs"},
	}
	var buf bytes.Buffer
	if err := writeHTML(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket.logs") {
		t.Error("expected added resource in output")
	}
	if !strings.Contains(out, "added") {
		t.Error("expected 'added' status class")
	}
}

func TestWriteHTML_Removed(t *testing.T) {
	result := diff.Result{
		Removed: []string{"aws_instance.web"},
	}
	var buf bytes.Buffer
	if err := writeHTML(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_instance.web") {
		t.Error("expected removed resource in output")
	}
	if !strings.Contains(out, "removed") {
		t.Error("expected 'removed' status")
	}
}

func TestWriteHTML_Changed(t *testing.T) {
	result := diff.Result{
		Changed: map[string]map[string]diff.AttrDiff{
			"aws_security_group.default": {
				"ingress": {Old: "0.0.0.0/0", New: "10.0.0.0/8"},
			},
		},
	}
	var buf bytes.Buffer
	if err := writeHTML(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_security_group.default") {
		t.Error("expected changed resource key")
	}
	if !strings.Contains(out, "ingress") {
		t.Error("expected attribute name")
	}
	if !strings.Contains(out, "10.0.0.0/8") {
		t.Error("expected new value")
	}
	if !strings.Contains(out, "0.0.0.0/0") {
		t.Error("expected old value")
	}
}

func TestWriteHTML_EscapesHTML(t *testing.T) {
	result := diff.Result{
		Added: []string{"aws_s3_bucket.<script>alert(1)</script>"},
	}
	var buf bytes.Buffer
	if err := writeHTML(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "<script>") {
		t.Error("expected HTML to be escaped")
	}
}
