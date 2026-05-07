package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftctl-diff configuration loaded from a YAML file.
type Config struct {
	// DefaultFormat is the output format used when --format is not supplied.
	DefaultFormat string `yaml:"default_format"`

	// DefaultOutput is the output destination used when --output is not supplied.
	DefaultOutput string `yaml:"default_output"`

	// IgnoreFile is the path to the ignore rules file.
	IgnoreFile string `yaml:"ignore_file"`

	// BaselineFile is the path to the baseline snapshot file.
	BaselineFile string `yaml:"baseline_file"`

	// IncludeTypes is a default list of resource types to include.
	IncludeTypes []string `yaml:"include_types"`

	// ExcludeTypes is a default list of resource types to exclude.
	ExcludeTypes []string `yaml:"exclude_types"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		DefaultFormat: "text",
		DefaultOutput: "-",
		IgnoreFile:    ".driftignore",
		BaselineFile:  ".drift-baseline.json",
	}
}

// Load reads a YAML config file from path and merges it on top of defaults.
// If the file does not exist, defaults are returned without error.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	return cfg, nil
}
