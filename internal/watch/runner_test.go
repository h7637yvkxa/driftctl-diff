package watch

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func writeStateFile(t *testing.T, serial int, resources []map[string]any) string {
	t.Helper()
	type tfstate struct {
		Version   int              `json:"version"`
		Serial    int              `json:"serial"`
		Resources []map[string]any `json:"resources"`
	}
	s := tfstate{Version: 4, Serial: serial, Resources: resources}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "*.tfstate")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(b)
	f.Close()
	return f.Name()
}

func TestRunOnce_NoDrift(t *testing.T) {
	res := []map[string]any{
		{"type": "aws_s3_bucket", "name": "main", "mode": "managed",
			"instances": []map[string]any{{"attributes": map[string]any{"id": "my-bucket"}}}},
	}
	base := writeStateFile(t, 1, res)
	target := writeStateFile(t, 2, res)

	var buf bytes.Buffer
	err := runOnce(RunOptions{BaseFile: base, TargetFile: target, Out: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' in output, got: %s", buf.String())
	}
}

func TestRunOnce_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := runOnce(RunOptions{
		BaseFile:   "/no/such/file.tfstate",
		TargetFile: "/no/such/file2.tfstate",
		Out:        &buf,
	})
	if err == nil {
		t.Fatal("expected error for missing state file")
	}
}

func TestRun_CancelImmediately(t *testing.T) {
	res := []map[string]any{}
	base := writeStateFile(t, 1, res)
	target := writeStateFile(t, 1, res)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	var buf bytes.Buffer
	err := Run(ctx, RunOptions{
		BaseFile:   base,
		TargetFile: target,
		Interval:   50 * time.Millisecond,
		Out:        &buf,
	})
	if err == nil {
		t.Fatal("expected context deadline error")
	}
	if !strings.Contains(buf.String(), "watching") {
		t.Errorf("expected startup line in output, got: %s", buf.String())
	}
}
