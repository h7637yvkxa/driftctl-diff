package lint

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds lint configuration loaded from a file.
type Config struct {
	DisabledRules []string `yaml:"disabled_rules"`
	FailOnSeverity string   `yaml:"fail_on_severity"` // "warning" | "error"
}

// DefaultLintConfig returns a sensible default lint config.
func DefaultLintConfig() Config {
	return Config{
		DisabledRules:  []string{},
		FailOnSeverity: "error",
	}
}

// LoadConfig reads a lint config YAML file from path.
// If the file does not exist, the default config is returned.
func LoadConfig(path string) (Config, error) {
	cfg := DefaultLintConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("lint: read config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("lint: parse config: %w", err)
	}
	return cfg, nil
}

// FilterRules returns only the rules that are not disabled by the config.
func FilterRules(rules []Rule, cfg Config) []Rule {
	disabled := make(map[string]bool, len(cfg.DisabledRules))
	for _, id := range cfg.DisabledRules {
		disabled[id] = true
	}
	var out []Rule
	for _, r := range rules {
		if !disabled[r.ID] {
			out = append(out, r)
		}
	}
	return out
}

// ShouldFail returns true if any finding meets or exceeds the configured severity threshold.
func ShouldFail(findings []Finding, cfg Config) bool {
	for _, f := range findings {
		if cfg.FailOnSeverity == "warning" {
			return true
		}
		if cfg.FailOnSeverity == "error" && f.Severity == "error" {
			return true
		}
	}
	return false
}
