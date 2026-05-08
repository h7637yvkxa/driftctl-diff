package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/owner/driftctl-diff/internal/diff"
)

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool    `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type sarifResult struct {
	RuleID  string        `json:"ruleId"`
	Level   string        `json:"level"`
	Message sarifMessage  `json:"message"`
	Locations []sarifLocation `json:"locations,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	LogicalLocations []sarifLogicalLocation `json:"logicalLocations"`
}

type sarifLogicalLocation struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func writeSARIF(w io.Writer, result diff.Result) error {
	var results []sarifResult

	keys := make([]string, 0, len(result.Added)+len(result.Removed)+len(result.Changed))
	for k := range result.Added {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		results = append(results, sarifResult{
			RuleID:  "drift/added",
			Level:   "warning",
			Message: sarifMessage{Text: fmt.Sprintf("Resource %q exists in source but not in target", k)},
			Locations: []sarifLocation{{LogicalLocations: []sarifLogicalLocation{{Name: k, Kind: "resource"}}}},
		})
	}

	keys = keys[:0]
	for k := range result.Removed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		results = append(results, sarifResult{
			RuleID:  "drift/removed",
			Level:   "warning",
			Message: sarifMessage{Text: fmt.Sprintf("Resource %q exists in target but not in source", k)},
			Locations: []sarifLocation{{LogicalLocations: []sarifLogicalLocation{{Name: k, Kind: "resource"}}}},
		})
	}

	keys = keys[:0]
	for k := range result.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		entry := result.Changed[k]
		for attr := range entry.Attributes {
			results = append(results, sarifResult{
				RuleID:  "drift/changed",
				Level:   "note",
				Message: sarifMessage{Text: fmt.Sprintf("Resource %q attribute %q differs between environments", k, attr)},
				Locations: []sarifLocation{{LogicalLocations: []sarifLogicalLocation{{Name: k, Kind: "resource"}}}},
			})
		}
	}

	log := sarifLog{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{{
			Tool:    sarifTool{Driver: sarifDriver{Name: "driftctl-diff", Version: "0.1.0"}},
			Results: results,
		}},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}
