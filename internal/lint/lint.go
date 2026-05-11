package lint

import (
	"fmt"
	"strings"

	"github.com/owner/driftctl-diff/internal/diff"
)

// Rule represents a lint rule that can be applied to a DriftResult.
type Rule struct {
	ID      string
	Message string
	Check   func(result diff.Result) []Finding
}

// Finding is a single lint violation.
type Finding struct {
	RuleID   string
	Severity string
	Resource string
	Detail   string
}

// String returns a human-readable representation of a Finding.
func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s (%s): %s", f.Severity, f.RuleID, f.Resource, f.Detail)
}

// DefaultRules returns the built-in lint rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			ID:      "DRIFT001",
			Message: "Resource added outside of Terraform",
			Check: func(r diff.Result) []Finding {
				var findings []Finding
				for key := range r.Added {
					findings = append(findings, Finding{
						RuleID:   "DRIFT001",
						Severity: "warning",
						Resource: key,
						Detail:   "resource exists in target but not in source state",
					})
				}
				return findings
			},
		},
		{
			ID:      "DRIFT002",
			Message: "Resource removed from state",
			Check: func(r diff.Result) []Finding {
				var findings []Finding
				for key := range r.Removed {
					findings = append(findings, Finding{
						RuleID:   "DRIFT002",
						Severity: "error",
						Resource: key,
						Detail:   "resource exists in source but is missing from target state",
					})
				}
				return findings
			},
		},
		{
			ID:      "DRIFT003",
			Message: "Resource attribute changed",
			Check: func(r diff.Result) []Finding {
				var findings []Finding
				for key, attrs := range r.Changed {
					keys := make([]string, 0, len(attrs))
					for attr := range attrs {
						keys = append(keys, attr)
					}
					findings = append(findings, Finding{
						RuleID:   "DRIFT003",
						Severity: "warning",
						Resource: key,
						Detail:   fmt.Sprintf("attributes changed: %s", strings.Join(keys, ", ")),
					})
				}
				return findings
			},
		},
	}
}

// Run applies all provided rules to the result and returns all findings.
func Run(result diff.Result, rules []Rule) []Finding {
	var all []Finding
	for _, rule := range rules {
		all = append(all, rule.Check(result)...)
	}
	return all
}
