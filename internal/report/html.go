package report

import (
	"fmt"
	"html"
	"io"
	"sort"

	"github.com/user/driftctl-diff/internal/diff"
)

func writeHTML(w io.Writer, result diff.Result) error {
	_, err := fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>driftctl-diff Report</title>
<style>
body{font-family:sans-serif;margin:2rem;background:#f9f9f9;color:#222}
h1{color:#333}table{border-collapse:collapse;width:100%;margin-top:1rem}
th{background:#4a90d9;color:#fff;padding:.5rem 1rem;text-align:left}
td{padding:.4rem 1rem;border-bottom:1px solid #ddd}
tr:nth-child(even){background:#eef2f7}
.added{color:#2a7a2a;font-weight:bold}
.removed{color:#a00;font-weight:bold}
.changed{color:#b36b00;font-weight:bold}
.none{color:#555}
</style>
</head>
<body>
<h1>Drift Report</h1>
`)
	if err != nil {
		return err
	}

	summaryClass := "none"
	if result.HasDrift() {
		summaryClass = "changed"
	}
	fmt.Fprintf(w, "<p class=\"%s\">%s</p>\n", summaryClass, html.EscapeString(result.Summary()))

	if !result.HasDrift() {
		_, err = fmt.Fprint(w, "</body>\n</html>\n")
		return err
	}

	fmt.Fprint(w, "<table>\n<thead><tr><th>Resource</th><th>Status</th><th>Attribute</th><th>Base</th><th>Current</th></tr></thead>\n<tbody>\n")

	writeHTMLRows(w, result.Added, "added", "added", "", "")
	writeHTMLRows(w, result.Removed, "removed", "removed", "", "")

	keys := make([]string, 0, len(result.Changed))
	for k := range result.Changed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		attrs := result.Changed[k]
		attrKeys := make([]string, 0, len(attrs))
		for a := range attrs {
			attrKeys = append(attrKeys, a)
		}
		sort.Strings(attrKeys)
		for _, a := range attrKeys {
			v := attrs[a]
			fmt.Fprintf(w, "<tr><td>%s</td><td class=\"changed\">changed</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
				html.EscapeString(k), html.EscapeString(a),
				html.EscapeString(fmt.Sprintf("%v", v.Old)),
				html.EscapeString(fmt.Sprintf("%v", v.New)))
		}
	}

	_, err = fmt.Fprint(w, "</tbody>\n</table>\n</body>\n</html>\n")
	return err
}

func writeHTMLRows(w io.Writer, keys []string, cssClass, status, valOld, valNew string) {
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	for _, k := range sorted {
		fmt.Fprintf(w, "<tr><td>%s</td><td class=\"%s\">%s</td><td>-</td><td>%s</td><td>%s</td></tr>\n",
			html.EscapeString(k), cssClass, status, valOld, valNew)
	}
}
