package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/owner/driftctl-diff/internal/diff"
)

// SonarQube Generic Issue Import Format
// https://docs.sonarqube.org/latest/analyzing-source-code/importing-external-issues/generic-issue-import-format/

type sonarReport struct {
	Issues []sonarIssue `json:"issues"`
}

type sonarIssue struct {
	EngineID        string         `json:"engineId"`
	RuleID          string         `json:"ruleId"`
	Severity        string         `json:"severity"`
	Type            string         `json:"type"`
	PrimaryLocation sonarLocation  `json:"primaryLocation"`
	SecondaryLocations []sonarLocation `json:"secondaryLocations,omitempty"`
}

type sonarLocation struct {
	Message   string `json:"message"`
	FilePath  string `json:"filePath"`
}

func writeSonarQube(w io.Writer, result diff.Result, stateFile string) error {
	if stateFile == "" {
		stateFile = "terraform.tfstate"
	}

	report := sonarReport{Issues: []sonarIssue{}}

	keys := make([]string, 0, len(result.Added)+len(result.Removed)+len(result.Changed))
	for k := range result.Added {
		keys = append(keys, "added:"+k)
	}
	for k := range result.Removed {
		keys = append(keys, "removed:"+k)
	}
	for k := range result.Changed {
		keys = append(keys, "changed:"+k)
	}
	sort.Strings(keys)

	for _, entry := range keys {
		var ruleID, severity, msg string
		switch entry[:entry.index(":")] {
		case "added":
			k := entry[6:]
			ruleID = "drift.resource.added"
			severity = "MAJOR"
			msg = fmt.Sprintf("Resource %q exists in target state but not in source", k)
		case "removed":
			k := entry[8:]
			ruleID = "drift.resource.removed"
			severity = "CRITICAL"
			msg = fmt.Sprintf("Resource %q exists in source state but is missing in target", k)
		case "changed":
			k := entry[8:]
			attrs := result.Changed[k]
			ruleID = "drift.resource.changed"
			severity = "MAJOR"
			msg = fmt.Sprintf("Resource %q has %d changed attribute(s): %v", k, len(attrs), attrKeys(attrs))
		}
		report.Issues = append(report.Issues, sonarIssue{
			EngineID: "driftctl-diff",
			RuleID:   ruleID,
			Severity: severity,
			Type:     "BUG",
			PrimaryLocation: sonarLocation{
				Message:  msg,
				FilePath: stateFile,
			},
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func (s string) index(sep string) int {
	return len(s) - len(s[len(sep):])
}

func attrKeys(attrs map[string]diff.AttributeDiff) []string {
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
