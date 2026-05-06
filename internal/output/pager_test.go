package output

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPagerWriter_FallbackOnNonTTY(t *testing.T) {
	// A bytes.Buffer is never a TTY, so the pager should not be invoked.
	var buf bytes.Buffer
	pw, cleanup, err := NewPagerWriter(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cleanup()

	msg := "hello from pager test"
	_, _ = pw.Write([]byte(msg))

	if !strings.Contains(buf.String(), msg) {
		t.Errorf("expected %q in output, got %q", msg, buf.String())
	}
}

func TestPagerWriter_WriteMultiple(t *testing.T) {
	var buf bytes.Buffer
	pw, cleanup, err := NewPagerWriter(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cleanup()

	lines := []string{"line one\n", "line two\n", "line three\n"}
	for _, l := range lines {
		if _, err := pw.Write([]byte(l)); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
	}

	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected %q in output", l)
		}
	}
}

func TestIsTTY_Buffer(t *testing.T) {
	var buf bytes.Buffer
	if isTTY(&buf) {
		t.Error("bytes.Buffer should not be detected as TTY")
	}
}

func TestIsTTY_File(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tty-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Regular files are not TTYs.
	if isTTY(f) {
		t.Error("regular file should not be detected as TTY")
	}
}

func TestDetectPager_EnvOverride(t *testing.T) {
	t.Setenv("PAGER", "cat")
	if got := detectPager(); got != "cat" {
		t.Errorf("expected 'cat', got %q", got)
	}
}

func TestDetectPager_EmptyEnv(t *testing.T) {
	t.Setenv("PAGER", "")
	// Should not panic; result depends on system binaries.
	_ = detectPager()
}
