// Package exitcode defines structured exit codes for driftctl-diff,
// allowing CI pipelines and scripts to distinguish between clean runs,
// detected drift, and hard errors.
package exitcode

import (
	"fmt"
	"os"

	"github.com/driftctl-diff/internal/diff"
)

// Code represents a process exit code.
type Code int

const (
	// OK indicates the comparison completed with no drift detected.
	OK Code = 0

	// DriftDetected indicates the comparison completed and drift was found.
	DriftDetected Code = 1

	// Error indicates a fatal error occurred (bad input, I/O failure, etc.).
	Error Code = 2
)

// String returns a human-readable label for the exit code.
func (c Code) String() string {
	switch c {
	case OK:
		return "ok"
	case DriftDetected:
		return "drift_detected"
	case Error:
		return "error"
	default:
		return fmt.Sprintf("unknown(%d)", int(c))
	}
}

// FromResult derives the appropriate exit code from a diff result.
// Returns DriftDetected if any added, removed, or changed resources exist;
// otherwise returns OK.
func FromResult(r *diff.Result) Code {
	if r == nil {
		return OK
	}
	if len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0 {
		return DriftDetected
	}
	return OK
}

// Exit calls os.Exit with the integer value of the code.
// This is a thin wrapper so callers can swap it out in tests.
var Exit = func(c Code) {
	os.Exit(int(c))
}
