package cmd

import (
	"fmt"
	"os"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/report"
	"github.com/owner/driftctl-diff/internal/state"
	"github.com/spf13/cobra"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook <base> <target> <url>",
	Short: "Compare state files and POST a drift report to a webhook URL",
	Args:  cobra.ExactArgs(3),
	RunE:  runWebhook,
}

func init() {
	rootCmd.AddCommand(webhookCmd)
}

func runWebhook(cmd *cobra.Command, args []string) error {
	baseFile, targetFile, webhookURL := args[0], args[1], args[2]

	baseResources, err := state.ParseStateFile(baseFile)
	if err != nil {
		return fmt.Errorf("parse base state: %w", err)
	}
	targetResources, err := state.ParseStateFile(targetFile)
	if err != nil {
		return fmt.Errorf("parse target state: %w", err)
	}

	baseIndex := state.IndexResources(baseResources)
	targetIndex := state.IndexResources(targetResources)

	result := diff.Compare(baseIndex, targetIndex)

	fmt := report.NewFormatter()
	if err := fmt.WriteWebhook(os.Stdout, result, webhookURL); err != nil {
		return fmt.Errorf("webhook delivery failed: %w", err)
	}
	return nil
}
