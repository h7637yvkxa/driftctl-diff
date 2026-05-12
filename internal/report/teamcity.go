package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftctl-diff/internal/diff"
)

// writeTeamCity emits TeamCity service messages for the given diff result.
// These messages are parsed by the TeamCity build agent and surfaced in the UI.
func writeTeamCity(w io.Writer, result diff.Result) error {
	if !result.HasDrift() {
		_, err := fmt.Fprintln(w, "##teamcity[message text='No drift detected' status='NORMAL']")
		return err
	}

	for _, r := range result.Added {
		msg := tcEscape(fmt.Sprintf("Resource added: %s", r.Address))
		if _, err := fmt.Fprintf(w, "##teamcity[message text='%s' status='WARNING']\n", msg); err != nil {
			return err
		}
	}

	for _, r := range result.Removed {
		msg := tcEscape(fmt.Sprintf("Resource removed: %s", r.Address))
		if _, err := fmt.Fprintf(w, "##teamcity[message text='%s' status='WARNING']\n", msg); err != nil {
			return err
		}
	}

	for _, r := range result.Changed {
		for attr := range r.Attributes {
			msg := tcEscape(fmt.Sprintf("Attribute drift: %s → %s", r.Address, attr))
			if _, err := fmt.Fprintf(w, "##teamcity[message text='%s' status='WARNING']\n", msg); err != nil {
				return err
			}
		}
	}

	summary := tcEscape(result.Summary())
	_, err := fmt.Fprintf(w, "##teamcity[buildProblem description='%s']\n", summary)
	return err
}

// tcEscape escapes special characters in TeamCity service message values.
func tcEscape(s string) string {
	replacer := strings.NewReplacer(
		"'", "|\'|",
		"|", "||",
		"[", "|[",
		"]", "|]",
		"\n", "|n",
		"\r", "|r",
	)
	return replacer.Replace(s)
}
