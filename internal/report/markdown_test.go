package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatter_Markdown_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, "markdown")

	if err := f.Write(noDriftResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got:\n%s", out)
	}
}

func TestFormatter_Markdown_Drift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, "markdown")

	if err := f.Write(driftResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "## Drift Report") {
		t.Errorf("expected markdown heading, got:\n%s", out)
	}
	if !strings.Contains(out, "Added Resources") {
		t.Errorf("expected Added section, got:\n%s", out)
	}
	if !strings.Contains(out, "Removed Resources") {
		t.Errorf("expected Removed section, got:\n%s", out)
	}
	if !strings.Contains(out, "Changed Resources") {
		t.Errorf("expected Changed section, got:\n%s", out)
	}
	if !strings.Contains(out, "| Attribute | Base | Target |") {
		t.Errorf("expected attribute table header, got:\n%s", out)
	}
}

func TestFormatter_Markdown_ChangedAttributes(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, "markdown")

	if err := f.Write(driftResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "instance_type") {
		t.Errorf("expected attribute name in output, got:\n%s", out)
	}
}
