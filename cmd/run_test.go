package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "state-*.tfstate")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const sampleState = `{"version":4,"terraform_version":"1.5.0","resources":[{"type":"aws_instance","name":"web","provider":"provider[\"registry.terraform.io/hashicorp/aws\"]","instances":[{"attributes":{"ami":"ami-123","instance_type":"t3.micro"}}]}]}`

func TestRunDiff_NoDrift(t *testing.T) {
	base := writeTempState(t, sampleState)
	target := writeTempState(t, sampleState)

	baseFile = base
	targetFile = target
	outputFmt = "text"
	outputFile = filepath.Join(t.TempDir(), "out.txt")

	err := runDiff(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(outputFile)
	if !strings.Contains(string(data), "No drift") {
		t.Errorf("expected 'No drift' in output, got:\n%s", data)
	}
}

func TestRunDiff_JSONOutput(t *testing.T) {
	base := writeTempState(t, sampleState)

	changed := `{"version":4,"terraform_version":"1.5.0","resources":[{"type":"aws_instance","name":"web","provider":"provider[\"registry.terraform.io/hashicorp/aws\"]","instances":[{"attributes":{"ami":"ami-456","instance_type":"t3.large"}}]}]}`
	target := writeTempState(t, changed)

	baseFile = base
	targetFile = target
	outputFmt = "json"
	outputFile = filepath.Join(t.TempDir(), "out.json")

	// runDiff exits 2 on drift; capture without os.Exit by checking error path
	_ = runDiff(rootCmd, nil)

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&result); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, data)
	}

	if _, ok := result["changed"]; !ok {
		t.Errorf("expected 'changed' key in JSON output")
	}
}
