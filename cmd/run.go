package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/driftctl-diff/internal/diff"
	"github.com/user/driftctl-diff/internal/report"
	"github.com/user/driftctl-diff/internal/state"
)

func runDiff(cmd *cobra.Command, args []string) error {
	baseState, err := state.ParseStateFile(baseFile)
	if err != nil {
		return fmt.Errorf("reading base state %q: %w", baseFile, err)
	}

	targetState, err := state.ParseStateFile(targetFile)
	if err != nil {
		return fmt.Errorf("reading target state %q: %w", targetFile, err)
	}

	baseIndex := state.IndexResources(baseState)
	targetIndex := state.IndexResources(targetState)

	result := diff.Compare(baseIndex, targetIndex)

	var out *os.File
	if outputFile != "" {
		out, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("creating output file %q: %w", outputFile, err)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Comparing %s → %s\n", baseFile, targetFile)

	formatter := report.NewFormatter(out, outputFmt)
	if err := formatter.Write(result); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	if result.HasDrift() {
		os.Exit(2)
	}
	return nil
}
