package report

import (
	"fmt"
	"io"

	"github.com/owner/driftctl-diff/internal/diff"
)

// Format is the output format identifier.
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
	FormatCSV      Format = "csv"
	FormatHTML     Format = "html"
	FormatSARIF    Format = "sarif"
)

// Formatter writes a diff.Result to an io.Writer in the requested format.
type Formatter struct {
	format Format
}

// NewFormatter creates a Formatter for the given format string.
// Returns an error if the format is not recognised.
func NewFormatter(format string) (*Formatter, error) {
	f := Format(format)
	switch f {
	case FormatText, FormatJSON, FormatMarkdown, FormatCSV, FormatHTML, FormatSARIF:
		return &Formatter{format: f}, nil
	default:
		return nil, fmt.Errorf("unsupported format %q; valid options: text, json, markdown, csv, html, sarif", format)
	}
}

// Write renders result to w using the configured format.
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
	default:
		return fmt.Errorf("unsupported format %q", f.format)
	}
}
