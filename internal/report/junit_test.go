package report

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/state"
)

func parseJUnit(t *testing.T, data string) junitTestSuites {
	t.Helper()
	var suites junitTestSuites
	if err := xml.Unmarshal([]byte(data), &suites); err != nil {
		t.Fatalf("failed to parse junit xml: %v", err)
	}
	return suites
}

func TestWriteJUnit_NoDrift(t *testing.T) {
	result := diff.Result{}
	var buf bytes.Buffer
	if err := writeJUnit(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	suites := parseJUnit(t, buf.String())
	if len(suites.Suites) != 1 {
		t.Fatalf("expected 1 suite, got %d", len(suites.Suites))
	}
	s := suites.Suites[0]
	if s.Failures != 0 {
		t.Errorf("expected 0 failures, got %d", s.Failures)
	}
	if s.Tests != 0 {
		t.Errorf("expected 0 tests, got %d", s.Tests)
	}
}

func TestWriteJUnit_Added(t *testing.T) {
	res := state.Resource{Type: "aws_s3_bucket", Name: "logs"}
	result := diff.Result{
		Added: map[string]state.Resource{"aws_s3_bucket.logs": res},
	}
	var buf bytes.Buffer
	if err := writeJUnit(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket.logs") {
		t.Errorf("expected resource key in output")
	}
	suites := parseJUnit(t, out)
	if suites.Suites[0].Failures != 1 {
		t.Errorf("expected 1 failure")
	}
}

func TestWriteJUnit_Removed(t *testing.T) {
	res := state.Resource{Type: "aws_instance", Name: "web"}
	result := diff.Result{
		Removed: map[string]state.Resource{"aws_instance.web": res},
	}
	var buf bytes.Buffer
	if err := writeJUnit(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	suites := parseJUnit(t, buf.String())
	if suites.Suites[0].Failures != 1 {
		t.Errorf("expected 1 failure")
	}
}

func TestWriteJUnit_Changed(t *testing.T) {
	src := state.Resource{Type: "aws_instance", Name: "app"}
	dst := state.Resource{Type: "aws_instance", Name: "app"}
	result := diff.Result{
		Changed: map[string]diff.Change{
			"aws_instance.app": {
				Source: src,
				Target: dst,
				Attributes: map[string]diff.AttributeDiff{
					"instance_type": {Source: "t2.micro", Target: "t3.small"},
				},
			},
		},
	}
	var buf bytes.Buffer
	if err := writeJUnit(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "1 attribute(s) changed") {
		t.Errorf("expected attribute count in failure message, got:\n%s", out)
	}
}
