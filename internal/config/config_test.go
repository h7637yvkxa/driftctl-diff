package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftctl-diff/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.DefaultFormat != "text" {
		t.Errorf("expected default_format=text, got %q", cfg.DefaultFormat)
	}
	if cfg.DefaultOutput != "-" {
		t.Errorf("expected default_output=-, got %q", cfg.DefaultOutput)
	}
	if cfg.IgnoreFile != ".driftignore" {
		t.Errorf("expected ignore_file=.driftignore, got %q", cfg.IgnoreFile)
	}
}

func TestLoad_Missing(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg.DefaultFormat != "text" {
		t.Errorf("expected defaults when file missing, got format=%q", cfg.DefaultFormat)
	}
}

func TestLoad_Valid(t *testing.T) {
	yaml := `
default_format: json
default_output: report.txt
ignore_file: custom.driftignore
baseline_file: snap.json
include_types:
  - aws_instance
  - aws_s3_bucket
exclude_types:
  - aws_iam_role
`
	p := writeTempConfig(t, yaml)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultFormat != "json" {
		t.Errorf("default_format: got %q, want json", cfg.DefaultFormat)
	}
	if cfg.DefaultOutput != "report.txt" {
		t.Errorf("default_output: got %q, want report.txt", cfg.DefaultOutput)
	}
	if len(cfg.IncludeTypes) != 2 {
		t.Errorf("include_types: got %d items, want 2", len(cfg.IncludeTypes))
	}
	if len(cfg.ExcludeTypes) != 1 || cfg.ExcludeTypes[0] != "aws_iam_role" {
		t.Errorf("exclude_types: unexpected value %v", cfg.ExcludeTypes)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeTempConfig(t, ": bad: yaml: [")
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}
