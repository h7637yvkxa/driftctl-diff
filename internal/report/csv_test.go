package report

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/user/driftctl-diff/internal/diff"
)

func TestWriteCSV_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}
	if err := writeCSV(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	records := parseCSV(t, buf.String())
	if len(records) != 1 {
		t.Fatalf("expected only header row, got %d rows", len(records))
	}
	expectHeader(t, records[0])
}

func TestWriteCSV_Added(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Added: map[diff.ResourceKey]diff.Resource{
			{Type: "aws_s3_bucket", Name: "logs"}: {},
		},
	}
	if err := writeCSV(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	records := parseCSV(t, buf.String())
	if len(records) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(records))
	}
	row := records[1]
	if row[0] != "aws_s3_bucket" || row[1] != "logs" || row[2] != "added" {
		t.Errorf("unexpected row: %v", row)
	}
}

func TestWriteCSV_Removed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Removed: map[diff.ResourceKey]diff.Resource{
			{Type: "aws_instance", Name: "web"}: {},
		},
	}
	if err := writeCSV(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	records := parseCSV(t, buf.String())
	if len(records) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(records))
	}
	if records[1][2] != "removed" {
		t.Errorf("expected change_type=removed, got %s", records[1][2])
	}
}

func TestWriteCSV_Changed(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Changed: map[diff.ResourceKey]map[string]diff.AttributeChange{
			{Type: "aws_instance", Name: "app"}: {
				"instance_type": {From: "t2.micro", To: "t3.small"},
			},
		},
	}
	if err := writeCSV(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	records := parseCSV(t, buf.String())
	if len(records) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(records))
	}
	row := records[1]
	if row[2] != "changed" || row[3] != "instance_type" || row[4] != "t2.micro" || row[5] != "t3.small" {
		t.Errorf("unexpected changed row: %v", row)
	}
}

func parseCSV(t *testing.T, s string) [][]string {
	t.Helper()
	r := csv.NewReader(strings.NewReader(s))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	return records
}

func expectHeader(t *testing.T, row []string) {
	t.Helper()
	expected := []string{"type", "name", "change_type", "attribute", "baseline_value", "current_value"}
	for i, col := range expected {
		if row[i] != col {
			t.Errorf("header[%d]: want %q got %q", i, col, row[i])
		}
	}
}
