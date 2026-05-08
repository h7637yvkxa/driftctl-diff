package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/filter"
	"github.com/owner/driftctl-diff/internal/metrics"
	"github.com/owner/driftctl-diff/internal/state"
	"github.com/spf13/cobra"
)

func init() {
	metricsCmd := &cobra.Command{
		Use:   "metrics <base-state> <target-state>",
		Short: "Print drift metrics in Prometheus text format",
		Args:  cobra.ExactArgs(2),
		RunE:  runMetrics,
	}
	registerFilterFlags(metricsCmd)
	metricsCmd.Flags().StringSliceP("label", "l", nil, "extra labels in key=value format")
	RootCmd.AddCommand(metricsCmd)
}

func runMetrics(cmd *cobra.Command, args []string) error {
	start := time.Now()

	baseFile, targetFile := args[0], args[1]

	baseState, err := state.ParseStateFile(baseFile)
	if err != nil {
		return fmt.Errorf("reading base state: %w", err)
	}
	targetState, err := state.ParseStateFile(targetFile)
	if err != nil {
		return fmt.Errorf("reading target state: %w", err)
	}

	baseIdx := state.IndexResources(baseState)
	targetIdx := state.IndexResources(targetState)

	result := diff.Compare(baseIdx, targetIdx)

	fopts := buildFilterOptions(cmd)
	filtered := filter.DriftResult(&result, fopts)

	end := time.Now()
	m := metrics.Collect(result, start, end, filtered)

	labelStrs, _ := cmd.Flags().GetStringSlice("label")
	labels, err := parseMetricLabels(labelStrs)
	if err != nil {
		return err
	}

	return metrics.WritePrometheus(os.Stdout, m, labels)
}

func parseMetricLabels(raw []string) ([]metrics.Label, error) {
	out := make([]metrics.Label, 0, len(raw))
	for _, s := range raw {
		var k, v string
		if _, err := fmt.Sscanf(s, "%s", &k); err != nil || k == "" {
			return nil, fmt.Errorf("invalid label %q: expected key=value", s)
		}
		for i, c := range s {
			if c == '=' {
				k = s[:i]
				v = s[i+1:]
				break
			}
		}
		if k == "" {
			return nil, fmt.Errorf("invalid label %q: expected key=value", s)
		}
		out = append(out, metrics.Label{Key: k, Value: v})
	}
	return out, nil
}
