package report

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func decodeCycloneDX(t *testing.T, buf *bytes.Buffer) cycloneDXReport {
	t.Helper()
	// strip XML declaration before decoding
	body := strings.TrimPrefix(buf.String(), xml.Header)
	var report cycloneDXReport
	if err := xml.Unmarshal([]byte(body), &report); err != nil {
		t.Fatalf("xml.Unmarshal: %v", err)
	}
	return report
}

func TestWriteCycloneDX_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}
	if err := writeCycloneDX(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	report := decodeCycloneDX(t, &buf)
	if len(report.Components) != 0 {
		t.Errorf("expected 0 components, got %d", len(report.Components))
	}
	if report.Metadata.Tools[0].Name != "driftctl-diff" {
		t.Errorf("unexpected tool name: %s", report.Metadata.Tools[0].Name)
	}
}

func TestWriteCycloneDX_Added(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Added: map[string]map[string]interface{}{
			"aws_s3_bucket.logs": {"bucket": "logs"},
		},
	}
	if err := writeCycloneDX(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	report := decodeCycloneDX(t, &buf)
	if len(report.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(report.Components))
	}
	c := report.Components[0]
	if c.Name != "aws_s3_bucket.logs" {
		t.Errorf("unexpected name: %s", c.Name)
	}
	if c.Properties[0].Value != "added" {
		t.Errorf("expected status=added, got %s", c.Properties[0].Value)
	}
}

func TestWriteCycloneDX_Removed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Removed: map[string]map[string]interface{}{
			"aws_instance.web": {"ami": "ami-123"},
		},
	}
	if err := writeCycloneDX(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	report := decodeCycloneDX(t, &buf)
	if len(report.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(report.Components))
	}
	if report.Components[0].Properties[0].Value != "removed" {
		t.Errorf("expected removed status")
	}
}

func TestWriteCycloneDX_Changed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Changed: map[string][]diff.AttributeChange{
			"aws_security_group.main": {
				{Attribute: "ingress", OldValue: "0.0.0.0/0", NewValue: "10.0.0.0/8"},
			},
		},
	}
	if err := writeCycloneDX(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	report := decodeCycloneDX(t, &buf)
	if len(report.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(report.Components))
	}
	props := report.Components[0].Properties
	if len(props) < 2 {
		t.Fatalf("expected at least 2 properties, got %d", len(props))
	}
	if props[0].Value != "changed" {
		t.Errorf("expected status=changed")
	}
	if !strings.Contains(props[1].Name, "ingress") {
		t.Errorf("expected ingress attribute property, got %s", props[1].Name)
	}
}
