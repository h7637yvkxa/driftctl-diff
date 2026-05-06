package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved snapshot of a diff result used as a reference point.
type Baseline struct {
	CreatedAt  time.Time          `json:"created_at"`
	SourceFile string             `json:"source_file"`
	TargetFile string             `json:"target_file"`
	Entries    []BaselineEntry    `json:"entries"`
}

// BaselineEntry records a known drift item that should be suppressed in future diffs.
type BaselineEntry struct {
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	Attribute    string `json:"attribute,omitempty"` // empty means entire resource
}

// Save writes the baseline to the given file path as JSON.
func Save(path string, b *Baseline) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// Load reads a baseline from the given file path.
func Load(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: file not found: %s", path)
		}
		return nil, fmt.Errorf("baseline: open %s: %w", path, err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode %s: %w", path, err)
	}
	return &b, nil
}

// EntryKey returns a comparable string key for a BaselineEntry.
func EntryKey(e BaselineEntry) string {
	if e.Attribute != "" {
		return e.ResourceType + "." + e.ResourceName + ":" + e.Attribute
	}
	return e.ResourceType + "." + e.ResourceName
}
