package watch

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/state"
)

// RunOptions configures a watch-mode run.
type RunOptions struct {
	BaseFile   string
	TargetFile string
	Interval   time.Duration
	Out        io.Writer
}

// Run starts watch mode: it re-runs the diff every time either state file
// changes and prints a compact summary to Out.
func Run(ctx context.Context, opts RunOptions) error {
	w := &FileWatcher{
		BaseFile:   opts.BaseFile,
		TargetFile: opts.TargetFile,
		Interval:   opts.Interval,
	}

	ch := make(chan ChangeEvent, 4)
	errCh := make(chan error, 1)
	go func() { errCh <- w.Watch(ctx, ch) }()

	fmt.Fprintf(opts.Out, "watching %s and %s (interval %s)\n",
		opts.BaseFile, opts.TargetFile, opts.Interval)

	// Run once immediately.
	if err := runOnce(opts); err != nil {
		fmt.Fprintf(opts.Out, "[error] %v\n", err)
	}

	for {
		select {
		case evt := <-ch:
			fmt.Fprintf(opts.Out, "[%s] change detected (base=%v target=%v) — re-running diff\n",
				evt.At.Format(time.RFC3339), evt.BaseChanged, evt.TargetChanged)
			if err := runOnce(opts); err != nil {
				fmt.Fprintf(opts.Out, "[error] %v\n", err)
			}
		case err := <-errCh:
			return err
		}
	}
}

func runOnce(opts RunOptions) error {
	baseIdx, err := loadIndex(opts.BaseFile)
	if err != nil {
		return fmt.Errorf("base: %w", err)
	}
	targetIdx, err := loadIndex(opts.TargetFile)
	if err != nil {
		return fmt.Errorf("target: %w", err)
	}

	result := diff.Compare(baseIdx, targetIdx)
	fmt.Fprintf(opts.Out, "  %s\n", result.Summary())
	return nil
}

func loadIndex(path string) (map[string]state.Resource, error) {
	sf, err := state.ParseStateFile(path)
	if err != nil {
		return nil, err
	}
	return state.IndexResources(sf), nil
}
