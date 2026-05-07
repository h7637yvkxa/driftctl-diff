package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/snapshot"
)

func makeResult() diff.Result {
	return diff.Result{
		Added:   []diff.Resource{{Key: "aws_s3_bucket.logs"}},
		Removed: []diff.Resource{{Key: "aws_instance.web"}},
		Changed: []diff.Change{{Key: "aws_db_instance.main"}},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	result := makeResult()
	if err := snapshot.Save(path, result, "test-label", "env/prod"); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if snap.Meta.Label != "test-label" {
		t.Errorf("label = %q, want %q", snap.Meta.Label, "test-label")
	}
	if snap.Meta.Source != "env/prod" {
		t.Errorf("source = %q, want %q", snap.Meta.Source, "env/prod")
	}
	if len(snap.Result.Added) != 1 || snap.Result.Added[0].Key != "aws_s3_bucket.logs" {
		t.Errorf("unexpected Added: %+v", snap.Result.Added)
	}
}

func TestLoad_Missing(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := snapshot.Save("/proc/invalid/nested/snap.json", makeResult(), "", "")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestSave_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "snap.json")
	if err := snapshot.Save(path, makeResult(), "lbl", "src"); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
