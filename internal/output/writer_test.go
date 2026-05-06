package output_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/output"
)

func TestNewDestination_Stdout(t *testing.T) {
	for _, path := range []string{"", "-"} {
		d, err := output.NewDestination(path)
		if err != nil {
			t.Fatalf("path=%q: unexpected error: %v", path, err)
		}
		if d.IsFile() {
			t.Errorf("path=%q: expected IsFile()=false", path)
		}
		if d.Writer() == nil {
			t.Errorf("path=%q: Writer() must not be nil", path)
		}
		// Closing stdout destination must not error.
		if err := d.Close(); err != nil {
			t.Errorf("path=%q: Close() error: %v", path, err)
		}
	}
}

func TestNewDestination_File(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "report.txt")

	d, err := output.NewDestination(outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.IsFile() {
		t.Error("expected IsFile()=true")
	}

	_, werr := strings.NewReader("hello drift").WriteTo(d.Writer())
	if werr != nil {
		t.Fatalf("write error: %v", werr)
	}

	if err := d.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}

	data, rerr := os.ReadFile(outPath)
	if rerr != nil {
		t.Fatalf("ReadFile error: %v", rerr)
	}
	if string(data) != "hello drift" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestNewDestination_InvalidPath(t *testing.T) {
	_, err := output.NewDestination("/nonexistent/dir/report.txt")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
