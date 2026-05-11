package lint

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempLintConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "lint.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp lint config: %v", err)
	}
	return p
}

func TestDefaultLintConfig(t *testing.T) {
	cfg := DefaultLintConfig()
	if cfg.FailOnSeverity != "error" {
		t.Errorf("expected 'error', got %q", cfg.FailOnSeverity)
	}
	if len(cfg.DisabledRules) != 0 {
		t.Errorf("expected no disabled rules by default")
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/lint.yaml")
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if cfg.FailOnSeverity != "error" {
		t.Errorf("expected default fail_on_severity")
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeTempLintConfig(t, "disabled_rules:\n  - DRIFT001\nfail_on_severity: warning\n")
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.DisabledRules) != 1 || cfg.DisabledRules[0] != "DRIFT001" {
		t.Errorf("unexpected disabled rules: %v", cfg.DisabledRules)
	}
	if cfg.FailOnSeverity != "warning" {
		t.Errorf("expected 'warning', got %q", cfg.FailOnSeverity)
	}
}

func TestFilterRules_DisableOne(t *testing.T) {
	cfg := Config{DisabledRules: []string{"DRIFT001"}}
	rules := FilterRules(DefaultRules(), cfg)
	for _, r := range rules {
		if r.ID == "DRIFT001" {
			t.Error("DRIFT001 should have been filtered out")
		}
	}
	if len(rules) != len(DefaultRules())-1 {
		t.Errorf("expected %d rules, got %d", len(DefaultRules())-1, len(rules))
	}
}

func TestShouldFail_ErrorSeverity(t *testing.T) {
	cfg := Config{FailOnSeverity: "error"}
	findings := []Finding{{Severity: "error"}}
	if !ShouldFail(findings, cfg) {
		t.Error("expected ShouldFail to return true for error finding")
	}
}

func TestShouldFail_NoFindings(t *testing.T) {
	cfg := Config{FailOnSeverity: "error"}
	if ShouldFail(nil, cfg) {
		t.Error("expected ShouldFail to return false with no findings")
	}
}
