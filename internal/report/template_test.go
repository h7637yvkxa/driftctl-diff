package report

import (
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func TestWriteTemplate_NoDrift(t *testing.T) {
	result := diff.Result{}
	var buf strings.Builder
	if err := writeTemplate(&buf, result, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No drift detected.") {
		t.Errorf("expected no-drift message, got:\n%s", out)
	}
}

func TestWriteTemplate_Drift(t *testing.T) {
	result := diff.Result{
		Added:   []diff.ResourceRef{{Type: "aws_s3_bucket", Name: "logs"}},
		Removed: []diff.ResourceRef{{Type: "aws_instance", Name: "web"}},
		Changed: map[string]map[string]diff.AttributeChange{
			"aws_security_group.default": {
				"description": {Before: "old", After: "new"},
			},
		},
	}
	var buf strings.Builder
	if err := writeTemplate(&buf, result, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{
		"aws_s3_bucket.logs",
		"aws_instance.web",
		"aws_security_group.default",
		"old → new",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWriteTemplate_CustomTemplate(t *testing.T) {
	result := diff.Result{
		Added: []diff.ResourceRef{{Type: "aws_vpc", Name: "main"}},
	}
	custom := `CUSTOM: added={{len .Added}}`
	var buf strings.Builder
	if err := writeTemplate(&buf, result, custom); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "CUSTOM: added=1" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestWriteTemplate_InvalidTemplate(t *testing.T) {
	result := diff.Result{}
	var buf strings.Builder
	err := writeTemplate(&buf, result, "{{.Unclosed")
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
	if !strings.Contains(err.Error(), "invalid template") {
		t.Errorf("unexpected error message: %v", err)
	}
}
