package watch

import (
	"context"
	"os"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestHashFile_Valid(t *testing.T) {
	path := writeTempFile(t, `{"version":4}`)
	h, err := hashFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h) != 64 {
		t.Fatalf("expected 64-char hex digest, got %d chars", len(h))
	}
}

func TestHashFile_Missing(t *testing.T) {
	_, err := hashFile("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestHashFile_Deterministic(t *testing.T) {
	path := writeTempFile(t, `{"version":4}`)
	h1, _ := hashFile(path)
	h2, _ := hashFile(path)
	if h1 != h2 {
		t.Fatal("hash should be deterministic")
	}
}

func TestWatch_EmitsOnChange(t *testing.T) {
	base := writeTempFile(t, `{"version":4,"serial":1}`)
	target := writeTempFile(t, `{"version":4,"serial":2}`)

	w := &FileWatcher{
		BaseFile:   base,
		TargetFile: target,
		Interval:   50 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := make(chan ChangeEvent, 1)
	go w.Watch(ctx, ch) //nolint:errcheck

	// Mutate the base file after the watcher has captured initial hashes.
	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(base, []byte(`{"version":4,"serial":99}`), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case evt := <-ch:
		if !evt.BaseChanged {
			t.Error("expected BaseChanged=true")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ChangeEvent")
	}
}

func TestWatch_CancelStops(t *testing.T) {
	base := writeTempFile(t, `{}`)
	target := writeTempFile(t, `{}`)

	w := &FileWatcher{BaseFile: base, TargetFile: target, Interval: 50 * time.Millisecond}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch := make(chan ChangeEvent, 1)
	err := w.Watch(ctx, ch)
	if err == nil {
		t.Fatal("expected context error after cancel")
	}
}
