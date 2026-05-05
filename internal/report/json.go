package report

import (
	"encoding/json"
	"io"

	"github.com/user/driftctl-diff/internal/diff"
)

// jsonReport is the serialisable representation of a drift result.
type jsonReport struct {
	HasDrift bool              `json:"has_drift"`
	Summary  jsonSummary       `json:"summary"`
	Added    []string          `json:"added"`
	Removed  []string          `json:"removed"`
	Changed  []jsonChange      `json:"changed"`
}

type jsonSummary struct {
	Added   int `json:"added"`
	Removed int `json:"removed"`
	Changed int `json:"changed"`
}

type jsonChange struct {
	Key        string              `json:"key"`
	Attributes map[string][2]string `json:"attributes"`
}

func writeJSON(w io.Writer, result *diff.Result) error {
	changes := make([]jsonChange, 0, len(result.Changed))
	for _, c := range result.Changed {
		changes = append(changes, jsonChange{
			Key:        c.Key,
			Attributes: c.Attributes,
		})
	}

	report := jsonReport{
		HasDrift: result.HasDrift(),
		Summary: jsonSummary{
			Added:   len(result.Added),
			Removed: len(result.Removed),
			Changed: len(result.Changed),
		},
		Added:   result.Added,
		Removed: result.Removed,
		Changed: changes,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
