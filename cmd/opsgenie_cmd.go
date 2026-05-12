package cmd

import (
	"fmt"
	"os"

	"github.com/owner/driftctl-diff/internal/report"
	"github.com/owner/driftctl-diff/internal/state"
	"github.com/spf13/cobra"
)

var opsgenieCmd = &cobra.Command{
	Use:   "opsgenie <base-state> <target-state>",
	Short: "Send a drift alert to OpsGenie",
	Args:  cobra.ExactArgs(2),
	RunE:  runOpsGenie,
}

func init() {
	opsgenieCmd.Flags().String("api-key", "", "OpsGenie API key (or set OPSGENIE_API_KEY)")
	opsgenieCmd.Flags().String("api-url", "", "OpsGenie API URL (default: https://api.opsgenie.com/v2/alerts)")
	registerFilterFlags(opsgenieCmd)
	Execute().AddCommand(opsgenieCmd)
}

func runOpsGenie(cmd *cobra.Command, args []string) error {
	apiKey, _ := cmd.Flags().GetString("api-key")
	if apiKey == "" {
		apiKey = os.Getenv("OPSGENIE_API_KEY")
	}
	if apiKey == "" {
		return fmt.Errorf("OpsGenie API key required (--api-key or OPSGENIE_API_KEY)")
	}
	apiURL, _ := cmd.Flags().GetString("api-url")

	baseState, err := state.ParseStateFile(args[0])
	if err != nil {
		return fmt.Errorf("parse base state: %w", err)
	}
	targetState, err := state.ParseStateFile(args[1])
	if err != nil {
		return fmt.Errorf("parse target state: %w", err)
	}

	baseIdx := state.IndexResources(baseState)
	targetIdx := state.IndexResources(targetState)

	result := diff.Compare(baseIdx, targetIdx)

	opts := buildFilterOptions(cmd)
	result = filter.DriftResult(result, opts)

	return report.WriteOpsGenie(os.Stdout, result, apiKey, apiURL)
}
