package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/driftctl-diff/internal/diff"
)

func writeMarkdown(w io.Writer, result *diff.Result) error {
	if !result.HasDrift() {
		_, err := fmt.Fprintln(w, "## Drift Report\n\n✅ **No drift detected.** All resources match across environments.")
		return err
	}

	var sb strings.Builder
	sb.WriteString("## Drift Report\n\n")
	sb.WriteString(fmt.Sprintf("> %s\n\n", result.Summary()))

	if len(result.Added) > 0 {
		sb.WriteString("### ➕ Added Resources\n\n")
		sb.WriteString("| Resource Key |\n")
		sb.WriteString("|---|\n")
		for _, key := range result.Added {
			sb.WriteString(fmt.Sprintf("| `%s` |\n", key))
		}
		sb.WriteString("\n")
	}

	if len(result.Removed) > 0 {
		sb.WriteString("### ➖ Removed Resources\n\n")
		sb.WriteString("| Resource Key |\n")
		sb.WriteString("|---|\n")
		for _, key := range result.Removed {
			sb.WriteString(fmt.Sprintf("| `%s` |\n", key))
		}
		sb.WriteString("\n")
	}

	if len(result.Changed) > 0 {
		sb.WriteString("### 🔄 Changed Resources\n\n")
		for _, change := range result.Changed {
			sb.WriteString(fmt.Sprintf("#### `%s`\n\n", change.Key))
			sb.WriteString("| Attribute | Base | Target |\n")
			sb.WriteString("|---|---|---|\n")
			for attr, delta := range change.Attributes {
				sb.WriteString(fmt.Sprintf("| `%s` | `%v` | `%v` |\n", attr, delta.Base, delta.Target))
			}
			sb.WriteString("\n")
		}
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}
