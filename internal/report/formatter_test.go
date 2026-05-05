package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/driftctl-diff/internal/diff"
	"github.com/user/driftctl-diff/internal/report"
)

func noDriftResult() *diff.Result {
	return &diff.Result{}
}

func driftResult() *diff.Result {
	return &diff.Result{
		Added:   []string{"aws_s3_bucket.new"},
		Removed: []string{"aws_s3_bucket.old"},
		Changed: []diff.Change{
			{
				Key: "aws_instance.web",
				Attributes: map[string][2]string{
					"instance_type": {"t2.micro", "t3.small"},
				},
			},
		},
	}
}

func TestFormatter_Text_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatText, &buf)
	if err := f.Write(noDriftResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestFormatter_Text_Drift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatText, &buf)
	if err := f.Write(driftResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Added", "Removed", "Changed", "instance_type", "t2.micro", "t3.small"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestFormatter_JSON_Drift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatJSON, &buf)
	if err := f.Write(driftResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["has_drift"] != true {
		t.Errorf("expected has_drift=true")
	}
	summary := out["summary"].(map[string]interface{})
	if summary["added"].(float64) != 1 {
		t.Errorf("expected 1 added")
	}
}
