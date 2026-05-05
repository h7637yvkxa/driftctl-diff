package state

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single resource entry from a Terraform state file.
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

// TerraformState represents the top-level structure of a Terraform state file.
type TerraformState struct {
	Version   int        `json:"version"`
	TFVersion string     `json:"terraform_version"`
	Resources []Resource `json:"resources"`
}

// ParseStateFile reads and parses a Terraform state file from the given path.
func ParseStateFile(path string) (*TerraformState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file %q: %w", path, err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file %q: %w", path, err)
	}

	return &state, nil
}

// ResourceKey returns a unique string key identifying a resource.
func ResourceKey(r Resource) string {
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

// IndexResources builds a map of resource key -> Resource for fast lookup.
func IndexResources(state *TerraformState) map[string]Resource {
	index := make(map[string]Resource, len(state.Resources))
	for _, r := range state.Resources {
		index[ResourceKey(r)] = r
	}
	return index
}
