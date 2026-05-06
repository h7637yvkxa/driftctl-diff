package output

import (
	"fmt"
	"io"
	"os"
)

// Destination represents where output should be written.
type Destination struct {
	w      io.WriteCloser
	isFile bool
}

// NewDestination returns a Destination writing to the given path.
// If path is empty or "-", output goes to stdout.
func NewDestination(path string) (*Destination, error) {
	if path == "" || path == "-" {
		return &Destination{w: os.Stdout, isFile: false}, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output: cannot open file %q: %w", path, err)
	}
	return &Destination{w: f, isFile: true}, nil
}

// Writer returns the underlying io.Writer.
func (d *Destination) Writer() io.Writer {
	return d.w
}

// Close closes the destination. A no-op when writing to stdout.
func (d *Destination) Close() error {
	if d.isFile {
		return d.w.Close()
	}
	return nil
}

// IsFile reports whether output is directed to a file.
func (d *Destination) IsFile() bool {
	return d.isFile
}
