package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func decodeSlack(t *testing.T, buf *bytes.Buffer) slackPayload {
	t.Helper()
	var p slackPayload
	if err := json.NewDecoder(buf).Decode(&p); err != nil {
		t.Fatalf("failed to decode slack payload: %v", err)
	}
	return p
}

func TestWriteSlack_NoDrift(t *testing.T) {
	result := diff.Result{}
	var buf bytes.Buffer
	if err := writeSlack(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := decodeSlack(t, &buf)
	if len(p.Blocks) == 0 {
		t.Fatal("expected at least one block")
	}
	header := p.Blocks[0]
	if header.Text == nil || !strings.Contains(header.Text.Text, "No drift") {
		t.Errorf("expected no-drift header, got: %+v", header.Text)
	}
}

func TestWriteSlack_Added(t *testing.T) {
	result := diff.Result{
		Added: map[string]map[string]interface{}{
			"aws_s3_bucket.logs": {"bucket": "logs"},
		},
	}
	var buf bytes.Buffer
	if err := writeSlack(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := decodeSlack(t, &buf)
	header := p.Blocks[0]
	if header.Text == nil || !strings.Contains(header.Text.Text, "drift detected") {
		t.Errorf("expected drift header, got: %+v", header.Text)
	}
	out := buf.String()
	_ = out // already decoded
	// check raw JSON contains resource key
	raw := new(bytes.Buffer)
	_ = writeSlack(raw, result)
	if !strings.Contains(raw.String(), "aws_s3_bucket.logs") {
		t.Error("expected resource key in slack output")
	}
}

func TestWriteSlack_Changed(t *testing.T) {
	result := diff.Result{
		Changed: map[string][]diff.AttributeDiff{
			"aws_instance.web": {
				{Attribute: "instance_type", OldValue: "t2.micro", NewValue: "t3.small"},
			},
		},
	}
	var buf bytes.Buffer
	if err := writeSlack(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "aws_instance.web") {
		t.Error("expected changed resource in slack output")
	}
	if !strings.Contains(buf.String(), "1 attribute") {
		t.Error("expected attribute count in slack output")
	}
}

func TestWriteSlack_SummaryFields(t *testing.T) {
	result := diff.Result{
		Added: map[string]map[string]interface{}{
			"aws_s3_bucket.a": {},
		},
		Removed: map[string]map[string]interface{}{
			"aws_s3_bucket.b": {},
		},
	}
	var buf bytes.Buffer
	if err := writeSlack(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p := decodeSlack(t, &buf)
	// block index 2 should be the summary fields block
	summaryBlock := p.Blocks[2]
	if len(summaryBlock.Fields) != 3 {
		t.Errorf("expected 3 summary fields, got %d", len(summaryBlock.Fields))
	}
}
