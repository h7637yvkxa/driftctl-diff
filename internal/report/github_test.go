package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func TestWriteGitHub_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}

	if err := writeGitHub(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "::notice") {
		t.Errorf("expected notice annotation, got: %s", out)
	}
	if strings.Contains(out, "::warning") || strings.Contains(out, "::error") {
		t.Errorf("expected no warning/error annotations for clean result, got: %s", out)
	}
}

func TestWriteGitHub_Added(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Added: map[string]map[string]interface{}{
			"aws_instance.web": {"ami": "ami-123"},
		},
	}

	if err := writeGitHub(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "::warning title=Drift Added::") {
		t.Errorf("expected warning annotation for added resource, got: %s", out)
	}
	if !strings.Contains(out, "aws_instance.web") {
		t.Errorf("expected resource key in output, got: %s", out)
	}
}

func TestWriteGitHub_Removed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Removed: map[string]map[string]interface{}{
			"aws_s3_bucket.logs": {"bucket": "my-logs"},
		},
	}

	if err := writeGitHub(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "::warning title=Drift Removed::") {
		t.Errorf("expected warning annotation for removed resource, got: %s", out)
	}
}

func TestWriteGitHub_Changed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Changed: map[string]map[string]diff.Delta{
			"aws_security_group.default": {
				"ingress": {Old: "0.0.0.0/0", New: "10.0.0.0/8"},
			},
		},
	}

	if err := writeGitHub(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "::error title=Drift Changed::") {
		t.Errorf("expected error annotation for changed resource, got: %s", out)
	}
	if !strings.Contains(out, "ingress") {
		t.Errorf("expected attribute name in output, got: %s", out)
	}
}
