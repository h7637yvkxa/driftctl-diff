package watch

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// FileWatcher polls a pair of state files and emits events when either changes.
type FileWatcher struct {
	BaseFile   string
	TargetFile string
	Interval   time.Duration

	hashes [2]string
}

// ChangeEvent is emitted when one or both watched files change.
type ChangeEvent struct {
	BaseChanged   bool
	TargetChanged bool
	At            time.Time
}

// Watch blocks until ctx is cancelled, sending a ChangeEvent on ch whenever
// either file's content hash changes.
func (w *FileWatcher) Watch(ctx context.Context, ch chan<- ChangeEvent) error {
	if w.Interval <= 0 {
		w.Interval = 5 * time.Second
	}

	// Capture initial hashes so we don't fire immediately.
	w.hashes[0], _ = hashFile(w.BaseFile)
	w.hashes[1], _ = hashFile(w.TargetFile)

	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			h0, _ := hashFile(w.BaseFile)
			h1, _ := hashFile(w.TargetFile)

			event := ChangeEvent{
				BaseChanged:   h0 != w.hashes[0],
				TargetChanged: h1 != w.hashes[1],
				At:            time.Now(),
			}

			if event.BaseChanged || event.TargetChanged {
				w.hashes[0] = h0
				w.hashes[1] = h1
				ch <- event
			}
		}
	}
}

// hashFile returns a hex SHA-256 digest of the file at path, or an empty
// string if the file cannot be read.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
