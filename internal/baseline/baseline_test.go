package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/owner/driftctl-diff/internal/baseline"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := &baseline.Baseline{
		CreatedAt:  time.Now().UTC().Truncate(time.Second),
		SourceFile: "env/prod.tfstate",
		TargetFile: "env/staging.tfstate",
		Entries: []baseline.BaselineEntry{
			{ResourceType: "aws_instance", ResourceName: "web"},
			{ResourceType: "aws_s3_bucket", ResourceName: "assets", Attribute: "tags"},
		},
	}

	if err := baseline.Save(path, b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.SourceFile != b.SourceFile {
		t.Errorf("SourceFile mismatch: got %q", loaded.SourceFile)
	}
}

func TestLoad_Missing(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestEntryKey(t *testing.T) {
	cases := []struct {
		entry baseline.BaselineEntry
		want  string
	}{
		{baseline.BaselineEntry{ResourceType: "aws_instance", ResourceName: "web"}, "aws_instance.web"},
		{baseline.BaselineEntry{ResourceType: "aws_s3_bucket", ResourceName: "logs", Attribute: "tags"}, "aws_s3_bucket.logs:tags"},
	}
	for _, tc := range cases {
		got := baseline.EntryKey(tc.entry)
		if got != tc.want {
			t.Errorf("EntryKey(%+v) = %q, want %q", tc.entry, got, tc.want)
		}
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := baseline.Save(filepath.Join(os.DevNull, "sub", "file.json"), &baseline.Baseline{})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
