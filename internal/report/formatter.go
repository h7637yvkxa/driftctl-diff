package report

import (
	"fmt"
	"io"

	"github.com/owner/driftctl-diff/internal/diff"
)

// Format enumerates supported output formats.
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
	FormatCSV      Format = "csv"
	FormatHTML     Format = "html"
	FormatSARIF    Format = "sarif"
	FormatJUnit    Format = "junit"
	FormatTemplate Format = "template"
	FormatSlack    Format = "slack"
	FormatGitLab   Format = "gitlab"
)

// Formatter writes a diff.Result to an io.Writer in a chosen format.
type Formatter struct {
	format       Format
	templatePath string
}

// NewFormatter returns a Formatter for the given format string.
func NewFormatter(format string, templatePath string) (*Formatter, error) {
	f := Format(format)
	switch f {
	case FormatText, FormatJSON, FormatMarkdown, FormatCSV,
		FormatHTML, FormatSARIF, FormatJUnit, FormatTemplate,
		FormatSlack, FormatGitLab:
		return &Formatter{format: f, templatePath: templatePath}, nil
	}
	return nil, fmt.Errorf("unsupported format: %q", format)
}

// Write renders result to w.
func (f *Formatter) Write(w io.Writer, result diff.Result) error {
	switch f.format {
	case FormatText:
		return writeText(w, result)
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
	}
	return fmt.Errorf("unsupported format: %q", f.format)
}
