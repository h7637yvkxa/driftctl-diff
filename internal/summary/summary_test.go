package summary_test

import (
	"testing"

	"github.com/your-org/driftctl-diff/internal/diff"
	"github.com/your-org/driftctl-diff/internal/summary"
)

func makeResource(rtype, name string) diff.Resource {
	return diff.Resource{Type: rtype, Name: name}
}

func TestBuild_Clean(t *testing.T) {
	result := diff.Result{}
	rep := summary.Build(result)
	if !rep.Clean {
		t.Fatal("expected Clean=true for empty result")
	}
	if rep.TotalDrift != 0 {
		t.Fatalf("expected TotalDrift=0, got %d", rep.TotalDrift)
	}
	if len(rep.ByType) != 0 {
		t.Fatalf("expected no type entries, got %d", len(rep.ByType))
	}
}

func TestBuild_Drift(t *testing.T) {
	result := diff.Result{
		Added: []diff.Resource{
			makeResource("aws_s3_bucket", "logs"),
			makeResource("aws_s3_bucket", "data"),
		},
		Removed: []diff.Resource{
			makeResource("aws_instance", "web"),
		},
		Changed: []diff.Resource{
			makeResource("aws_s3_bucket", "primary"),
			makeResource("aws_instance", "worker"),
		},
	}

	rep := summary.Build(result)

	if rep.Clean {
		t.Fatal("expected Clean=false")
	}
	if rep.TotalAdded != 2 {
		t.Fatalf("expected TotalAdded=2, got %d", rep.TotalAdded)
	}
	if rep.TotalRemoved != 1 {
		t.Fatalf("expected TotalRemoved=1, got %d", rep.TotalRemoved)
	}
	if rep.TotalChanged != 2 {
		t.Fatalf("expected TotalChanged=2, got %d", rep.TotalChanged)
	}
	if rep.TotalDrift != 5 {
		t.Fatalf("expected TotalDrift=5, got %d", rep.TotalDrift)
	}

	// Expect two type entries: aws_instance and aws_s3_bucket (sorted)
	if len(rep.ByType) != 2 {
		t.Fatalf("expected 2 type entries, got %d", len(rep.ByType))
	}
	if rep.ByType[0].Type != "aws_instance" {
		t.Fatalf("expected first type=aws_instance, got %s", rep.ByType[0].Type)
	}
	if rep.ByType[0].Removed != 1 || rep.ByType[0].Changed != 1 {
		t.Fatalf("aws_instance counts wrong: %+v", rep.ByType[0])
	}
	if rep.ByType[1].Type != "aws_s3_bucket" {
		t.Fatalf("expected second type=aws_s3_bucket, got %s", rep.ByType[1].Type)
	}
	if rep.ByType[1].Added != 2 || rep.ByType[1].Changed != 1 {
		t.Fatalf("aws_s3_bucket counts wrong: %+v", rep.ByType[1])
	}
}

func TestFormatLine(t *testing.T) {
	ts := summary.TypeSummary{Type: "aws_s3_bucket", Added: 2, Removed: 0, Changed: 1}
	line := summary.FormatLine(ts)
	if line == "" {
		t.Fatal("expected non-empty format line")
	}
}
