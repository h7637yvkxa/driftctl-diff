package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftctl-diff/internal/diff"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes a drift report to a writer.
type Formatter struct {
	Format Format
	Writer io.Writer
}

// NewFormatter creates a Formatter with the given format and writer.
func NewFormatter(format Format, w io.Writer) *Formatter {
	return &Formatter{Format: format, Writer: w}
}

// Write outputs the drift result using the configured format.
func (f *Formatter) Write(result *diff.Result) error {
	switch f.Format {
	case FormatJSON:
		return writeJSON(f.Writer, result)
	default:
		return writeText(f.Writer, result)
	}
}

func writeText(w io.Writer, result *diff.Result) error {
	if !result.HasDrift() {
		_, err := fmt.Fprintln(w, "✓ No drift detected.")
		return err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Drift detected: %d added, %d removed, %d changed\n",
		len(result.Added), len(result.Removed), len(result.Changed)))

	if len(result.Added) > 0 {
		sb.WriteString("\n[+] Added resources:\n")
		for _, r := range result.Added {
			sb.WriteString(fmt.Sprintf("  + %s\n", r))
		}
	}
	if len(result.Removed) > 0 {
		sb.WriteString("\n[-] Removed resources:\n")
		for _, r := range result.Removed {
			sb.WriteString(fmt.Sprintf("  - %s\n", r))
		}
	}
	if len(result.Changed) > 0 {
		sb.WriteString("\n[~] Changed resources:\n")
		for _, c := range result.Changed {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", c.Key))
			for attr, vals := range c.Attributes {
				sb.WriteString(fmt.Sprintf("      %s: %q -> %q\n", attr, vals[0], vals[1]))
			}
		}
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}
