package report

import (
	"fmt"
	"io"

	"github.com/your-org/driftctl-diff/internal/diff"
	"github.com/your-org/driftctl-diff/internal/summary"
)

// writeSummaryTable writes a per-type drift breakdown table to w.
func writeSummaryTable(w io.Writer, result diff.Result) {
	rep := summary.Build(result)

	if rep.Clean {
		fmt.Fprintln(w, "No drift detected — all resources are in sync.")
		return
	}

	fmt.Fprintln(w, "\nDrift Summary by Resource Type:")
	fmt.Fprintln(w, "----------------------------------------------------------------------")
	for _, ts := range rep.ByType {
		fmt.Fprintln(w, summary.FormatLine(ts))
	}
	fmt.Fprintln(w, "----------------------------------------------------------------------")
	fmt.Fprintf(w, "Total: %d added, %d removed, %d changed (%d drifted)\n",
		rep.TotalAdded, rep.TotalRemoved, rep.TotalChanged, rep.TotalDrift)
}
