package report

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/owner/driftctl-diff/internal/diff"
)

const defaultTemplate = `Drift Report
============
{{- if .NoDrift}}
No drift detected.
{{- else}}

Added Resources ({{len .Added}}):
{{- range .Added}}
  + {{.Type}}.{{.Name}}
{{- end}}

Removed Resources ({{len .Removed}}):
{{- range .Removed}}
  - {{.Type}}.{{.Name}}
{{- end}}

Changed Resources ({{len .Changed}}):
{{- range $key, $attrs := .Changed}}
  ~ {{$key}}:
{{- range $attr, $change := $attrs}}
    • {{$attr}}: {{$change.Before}} → {{$change.After}}
{{- end}}
{{- end}}
{{- end}}

Summary: {{.Summary}}
`

func writeTemplate(w io.Writer, result diff.Result, tmplStr string) error {
	if strings.TrimSpace(tmplStr) == "" {
		tmplStr = defaultTemplate
	}

	tmpl, err := template.New("drift").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	data := struct {
		NoDrift bool
		Added   []diff.ResourceRef
		Removed []diff.ResourceRef
		Changed map[string]map[string]diff.AttributeChange
		Summary string
	}{
		NoDrift: result.NoDrift(),
		Added:   result.Added,
		Removed: result.Removed,
		Changed: result.Changed,
		Summary: result.Summary(),
	}

	return tmpl.Execute(w, data)
}
