package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"

	"github.com/user/driftctl-diff/internal/diff"
)

// writeCSV writes the drift result as a CSV report to w.
// Columns: type, name, change_type, attribute, baseline_value, current_value
func writeCSV(w io.Writer, result diff.Result) error {
	cw := csv.NewWriter(w)

	header := []string{"type", "name", "change_type", "attribute", "baseline_value", "current_value"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("csv: write header: %w", err)
	}

	for _, r := range sortedKeys(result.Added) {
		if err := cw.Write([]string{r.Type, r.Name, "added", "", "", ""}); err != nil {
			return err
		}
	}

	for _, r := range sortedKeys(result.Removed) {
		if err := cw.Write([]string{r.Type, r.Name, "removed", "", "", ""}); err != nil {
			return err
		}
	}

	type changedRow struct {
		rType, rName, attr, baseline, current string
	}
	var rows []changedRow
	for key, attrs := range result.Changed {
		for attr, chg := range attrs {
			rows = append(rows, changedRow{
				rType:    key.Type,
				rName:    key.Name,
				attr:     attr,
				baseline: fmt.Sprintf("%v", chg.From),
				current:  fmt.Sprintf("%v", chg.To),
			})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].rType != rows[j].rType {
			return rows[i].rType < rows[j].rType
		}
		if rows[i].rName != rows[j].rName {
			return rows[i].rName < rows[j].rName
		}
		return rows[i].attr < rows[j].attr
	})
	for _, row := range rows {
		if err := cw.Write([]string{row.rType, row.rName, "changed", row.attr, row.baseline, row.current}); err != nil {
			return err
		}
	}

	cw.Flush()
	return cw.Error()
}

// sortedKeys returns resource keys from a map in deterministic order.
func sortedKeys[V any](m map[diff.ResourceKey]V) []diff.ResourceKey {
	keys := make([]diff.ResourceKey, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Type != keys[j].Type {
			return keys[i].Type < keys[j].Type
		}
		return keys[i].Name < keys[j].Name
	})
	return keys
}
