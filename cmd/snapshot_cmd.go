package cmd

import (
	"fmt"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/snapshot"
	"github.com/owner/driftctl-diff/internal/state"
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage diff snapshots for drift trend tracking",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save <base> <updated> <output>",
	Short: "Save current diff result as a snapshot",
	Args:  cobra.ExactArgs(3),
	RunE:  runSnapshotSave,
}

var snapshotCompareCmd = &cobra.Command{
	Use:   "compare <snapshot> <base> <updated>",
	Short: "Compare a saved snapshot against the current diff",
	Args:  cobra.ExactArgs(3),
	RunE:  runSnapshotCompare,
}

var snapshotLabel string

func init() {
	snapshotSaveCmd.Flags().StringVar(&snapshotLabel, "label", "", "Optional label for the snapshot")
	snapshotCmd.AddCommand(snapshotSaveCmd, snapshotCompareCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotSave(cmd *cobra.Command, args []string) error {
	result, err := buildDiffResult(args[0], args[1])
	if err != nil {
		return err
	}
	return snapshot.Save(args[2], result, snapshotLabel, args[0])
}

func runSnapshotCompare(cmd *cobra.Command, args []string) error {
	snap, err := snapshot.Load(args[0])
	if err != nil {
		return err
	}
	result, err := buildDiffResult(args[1], args[2])
	if err != nil {
		return err
	}
	delta := snap.CompareTo(result)
	fmt.Fprintln(cmd.OutOrStdout(), delta.Summary())
	return nil
}

func buildDiffResult(basePath, updatedPath string) (diff.Result, error) {
	base, err := state.ParseStateFile(basePath)
	if err != nil {
		return diff.Result{}, fmt.Errorf("base state: %w", err)
	}
	updated, err := state.ParseStateFile(updatedPath)
	if err != nil {
		return diff.Result{}, fmt.Errorf("updated state: %w", err)
	}
	return diff.Compare(state.IndexResources(base), state.IndexResources(updated)), nil
}
