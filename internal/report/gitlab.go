package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/owner/driftctl-diff/internal/diff"
)

// GitLab Code Quality report format
// https://docs.gitlab.com/ee/ci/testing/code_quality.html

type gitLabIssue struct {
	Description string          `json:"description"`
	Fingerprint string          `json:"fingerprint"`
	Severity    string          `json:"severity"`
	Location    gitLabLocation  `json:"location"`
}

type gitLabLocation struct {
	Path  string `json:"path"`
	Lines gitLabLines `json:"lines"`
}

type gitLabLines struct {
	Begin int `json:"begin"`
}

func writeGitLab(w io.Writer, result diff.Result) error {
	issues := []gitLabIssue{}

	keys := make([]string, 0, len(result.Added)+len(result.Removed)+len(result.Changed))
	for k := range result.Added {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		issues = append(issues, gitLabIssue{
			Description: fmt.Sprintf("Resource added: %s", k),
			Fingerprint: fmt.Sprintf("added-%s", k),
			Severity:    "minor",
			Location:    gitLabLocation{Path: "terraform.tfstate", Lines: gitLabLines{Begin: 1}},
		})
	}

	keys = keys[:0]
	for k := range result.Removed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		issues = append(issues, gitLabIssue{
			Description: fmt.Sprintf("Resource removed: %s", k),
			Fingerprint: fmt.Sprintf("removed-%s", k),
			Severity:    "major",
			Location:    gitLabLocation{Path: "terraform.tfstate", Lines: gitLabLines{Begin: 1}},
		})
	}

	keys = keys[:0]
	for k := range result.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		c := result.Changed[k]
		for attr := range c.Attributes {
			issues = append(issues, gitLabIssue{
				Description: fmt.Sprintf("Attribute drift in %s: %s", k, attr),
				Fingerprint: fmt.Sprintf("changed-%s-%s", k, attr),
				Severity:    "minor",
				Location:    gitLabLocation{Path: "terraform.tfstate", Lines: gitLabLines{Begin: 1}},
			})
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(issues)
}
