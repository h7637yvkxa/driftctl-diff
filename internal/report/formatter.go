package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/owner/driftctl-diff/internal/diff"
)

// Format enumerates supported output formats.
type Format string

const (
	FormatText      Format = "text"
	FormatJSON      Format = "json"
	FormatMarkdown  Format = "markdown"
	FormatCSV       Format = "csv"
	FormatHTML      Format = "html"
	FormatSARIF     Format = "sarif"
	FormatJUnit     Format = "junit"
	FormatTemplate  Format = "template"
	FormatSlack     Format = "slack"
	FormatGitLab    Format = "gitlab"
	FormatCycloneDX Format = "cyclonedx"
)

// Formatter writes a diff.Result in the requested format.
type Formatter struct {
	format       Format
	templatePath string
}

// NewFormatter constructs a Formatter for the given format string.
func NewFormatter(format string, templatePath string) (*Formatter, error) {
	f := Format(strings.ToLower(strings.TrimSpace(format)))
	switch f {
	case FormatText, FormatJSON, FormatMarkdown, FormatCSV, FormatHTML,
		FormatSARIF, FormatJUnit, FormatTemplate, FormatSlack, FormatGitLab,
		FormatCycloneDX:
		return &Formatter{format: f, templatePath: templatePath}, nil
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
}

// Write renders result to w using the configured format.
func (f *Formatter) Write(w io.Writer, result diff.Result) error {
	switch f.format {
	case FormatJSON:
		return writeJSON(w, result)
	case FormatMarkdown:
		return writeMarkdown(w, result)
	case FormatCSV:
		return writeCSV(w, result)
	case FormatHTML:
		return writeHTML(w, result)
	case FormatSARIF:
		return writeSARIF(w, result)
	case FormatJUnit:
		return writeJUnit(w, result)
	case FormatTemplate:
		return writeTemplate(w, result, f.templatePath)
	case FormatSlack:
		return writeSlack(w, result)
	case FormatGitLab:
		return writeGitLab(w, result)
	case FormatCycloneDX:
		return writeCycloneDX(w, result)
	default:
		return writeText(w, result)
	}
}

func writeText(w io.Writer, result diff.Result) error {
	if !result.HasDrift() {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}
	for key := range result.Added {
		fmt.Fprintf(w, "[+] %s\n", key)
	}
	for key := range result.Removed {
		fmt.Fprintf(w, "[-] %s\n", key)
	}
	for key, changes := range result.Changed {
		for _, ch := range changes {
			fmt.Fprintf(w, "[~] %s  %s: %v -> %v\n", key, ch.Attribute, ch.OldValue, ch.NewValue)
		}
	}
	return nil
}
