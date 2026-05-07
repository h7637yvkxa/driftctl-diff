package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
)

// Metadata holds information about when a snapshot was taken.
type Metadata struct {
	CreatedAt time.Time `json:"created_at"`
	Label     string    `json:"label"`
	Source    string    `json:"source"`
}

// Snapshot captures a diff result along with metadata for later comparison.
type Snapshot struct {
	Meta   Metadata    `json:"meta"`
	Result diff.Result `json:"result"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, result diff.Result, label, source string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}
	s := Snapshot{
		Meta: Metadata{
			CreatedAt: time.Now().UTC(),
			Label:     label,
			Source:    source,
		},
		Result: result,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: file not found: %s", path)
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}
