package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftctl-diff/internal/diff"
)

func TestWriteTeamCity_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}

	if err := writeTeamCity(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
	if strings.Contains(out, "buildProblem") {
		t.Errorf("expected no buildProblem when no drift, got: %s", out)
	}
}

func TestWriteTeamCity_Added(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Added: []diff.Resource{{Address: "aws_s3_bucket.logs"}},
	}

	if err := writeTeamCity(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket.logs") {
		t.Errorf("expected added resource in output, got: %s", out)
	}
	if !strings.Contains(out, "status='WARNING'") {
		t.Errorf("expected WARNING status, got: %s", out)
	}
	if !strings.Contains(out, "buildProblem") {
		t.Errorf("expected buildProblem when drift present, got: %s", out)
	}
}

func TestWriteTeamCity_Changed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Changed: []diff.ResourceDiff{
			{
				Address: "aws_instance.web",
				Attributes: map[string]diff.AttributeDiff{
					"instance_type": {Before: "t2.micro", After: "t3.small"},
				},
			},
		},
	}

	if err := writeTeamCity(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "aws_instance.web") {
		t.Errorf("expected changed resource in output, got: %s", out)
	}
	if !strings.Contains(out, "instance_type") {
		t.Errorf("expected attribute name in output, got: %s", out)
	}
}

func TestTCEscape(t *testing.T) {
	cases := []struct {
		input string
		want string
	}{
		{"hello", "hello"},
		{"say 'hi'", "say |\'|hi|\'|"},
		{"pipe|char", "pipe||char"},
		{"new\nline", "new|nline"},
	}
	for _, tc := range cases {
		got := tcEscape(tc.input)
		if got != tc.want {
			t.Errorf("tcEscape(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
