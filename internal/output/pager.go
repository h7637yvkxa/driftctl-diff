package output

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

// PagerWriter wraps a writer that pipes output through a terminal pager
// (e.g. less) when stdout is a TTY and a pager is available.
type PagerWriter struct {
	cmd    *exec.Cmd
	pipe   io.WriteCloser
	fallback io.Writer
}

// NewPagerWriter returns a PagerWriter if a pager binary is found and stdout
// is a terminal; otherwise it returns the provided fallback writer unchanged.
func NewPagerWriter(fallback io.Writer) (*PagerWriter, func(), error) {
	pager := detectPager()
	if pager == "" || !isTTY(fallback) {
		return &PagerWriter{fallback: fallback}, func() {}, nil
	}

	parts := strings.Fields(pager)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = fallback
	cmd.Stderr = os.Stderr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return &PagerWriter{fallback: fallback}, func() {}, nil
	}

	if err := cmd.Start(); err != nil {
		_ = pipe.Close()
		return &PagerWriter{fallback: fallback}, func() {}, nil
	}

	cleanup := func() {
		_ = pipe.Close()
		_ = cmd.Wait()
	}

	return &PagerWriter{cmd: cmd, pipe: pipe}, cleanup, nil
}

// Write implements io.Writer.
func (p *PagerWriter) Write(b []byte) (int, error) {
	if p.pipe != nil {
		return p.pipe.Write(b)
	}
	return p.fallback.Write(b)
}

// detectPager returns the pager command to use, checking PAGER env var first,
// then falling back to "less -FRX", then "more".
func detectPager() string {
	if v := os.Getenv("PAGER"); v != "" {
		return v
	}
	if _, err := exec.LookPath("less"); err == nil {
		return "less -FRX"
	}
	if _, err := exec.LookPath("more"); err == nil {
		return "more"
	}
	return ""
}

// isTTY reports whether w is an *os.File backed by a terminal.
func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
