package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/owner/driftctl-diff/internal/watch"
	"github.com/spf13/cobra"
)

var (
	watchInterval time.Duration
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <base-state> <target-state>",
		Short: "Re-run diff whenever a state file changes",
		Long: `Watch mode polls both state files at a configurable interval.
Whenever a change is detected the diff is re-executed and a summary
line is printed to stdout.`,
		Args: cobra.ExactArgs(2),
		RunE: runWatch,
	}

	watchCmd.Flags().DurationVarP(&watchInterval, "interval", "i", 10*time.Second,
		"polling interval (e.g. 5s, 1m)")

	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	base := args[0]
	target := args[1]

	if _, err := os.Stat(base); err != nil {
		return fmt.Errorf("base state file not found: %w", err)
	}
	if _, err := os.Stat(target); err != nil {
		return fmt.Errorf("target state file not found: %w", err)
	}

	ctx := cmd.Context()

	opts := watch.RunOptions{
		BaseFile:   base,
		TargetFile: target,
		Interval:   watchInterval,
		Out:        cmd.OutOrStdout(),
	}

	if err := watch.Run(ctx, opts); err != nil {
		// context.Canceled / DeadlineExceeded are expected on clean exit.
		if ctx.Err() != nil {
			return nil
		}
		return err
	}
	return nil
}
