package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/owner/driftctl-diff/internal/diff"
)

// githubAnnotation represents a GitHub Actions workflow command annotation.
type githubAnnotation struct {
	Level   string `json:"level"`
	File    string `json:"file"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// writeGitHub emits GitHub Actions workflow commands so that drift surfaces
// as annotations in pull-request checks, plus a machine-readable JSON array
// of all annotations written to w.
func writeGitHub(w io.Writer, result diff.Result) error {
	annotations := make([]githubAnnotation, 0)

	addedKeys := sortedKeys(result.Added)
	for _, k := range addedKeys {
		msg := fmt.Sprintf("Resource %q is present in source but missing in target", k)
		fmt.Fprintf(w, "::warning title=Drift Added::%s\n", msg)
		annotations = append(annotations, githubAnnotation{
			Level:   "warning",
			File:    "",
			Title:   "Drift Added",
			Message: msg,
		})
	}

	removedKeys := sortedKeys(result.Removed)
	for _, k := range removedKeys {
		msg := fmt.Sprintf("Resource %q exists in target but is absent from source", k)
		fmt.Fprintf(w, "::warning title=Drift Removed::%s\n", msg)
		annotations = append(annotations, githubAnnotation{
			Level:   "warning",
			File:    "",
			Title:   "Drift Removed",
			Message: msg,
		})
	}

	changedKeys := make([]string, 0, len(result.Changed))
	for k := range result.Changed {
		changedKeys = append(changedKeys, k)
	}
	sort.Strings(changedKeys)

	for _, k := range changedKeys {
		attrs := result.Changed[k]
		for attr, delta := range attrs {
			msg := fmt.Sprintf("Resource %q attribute %q changed: %v -> %v", k, attr, delta.Old, delta.New)
			fmt.Fprintf(w, "::error title=Drift Changed::%s\n", msg)
			annotations = append(annotations, githubAnnotation{
				Level:   "error",
				File:    "",
				Title:   "Drift Changed",
				Message: msg,
			})
		}
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "::notice title=Drift::No drift detected")
		return nil
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(annotations)
}
