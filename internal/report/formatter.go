package report

import (
	"fmt"
	"io"

	"github.com/user/driftctl-diff/internal/diff"
)

// Format enumerates supported output formats.
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
	FormatCSV      Format = "csv"
	FormatHTML     Format = "html"
)

// Formatter writes a diff result in the requested format.
type Formatter struct {
	format Format
}

// NewFormatter creates a Formatter for the given format string.
// An empty string defaults to text.
func NewFormatter(format string) (*Formatter, error) {
	if format == "" {
		format = string(FormatText)
	}
	switch Format(format) {
	case FormatText, FormatJSON, FormatMarkdown, FormatCSV, FormatHTML:
		return &Formatter{format: Format(format)}, nil
	default:
		return nil, fmt.Errorf("unsupported format %q; choose text, json, markdown, csv, or html", format)
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
	default:
		return writeText(w, result)
	}
}
