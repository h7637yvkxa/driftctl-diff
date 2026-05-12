package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/driftctl-diff/internal/filter"
)

// filterFlagSet holds raw CLI flag values for resource filtering.
type filterFlagSet struct {
	includeTypes string
	excludeTypes string
	includeNames string
	excludeNames string
}

var filterFlags filterFlagSet

// registerFilterFlags attaches filter-related flags to the given command.
func registerFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&filterFlags.includeTypes, "include-types", "",
		"Comma-separated list of resource types to include (e.g. aws_instance,aws_s3_bucket)")
	cmd.Flags().StringVar(&filterFlags.excludeTypes, "exclude-types", "",
		"Comma-separated list of resource types to exclude")
	cmd.Flags().StringVar(&filterFlags.includeNames, "include-names", "",
		"Comma-separated list of resource names to include")
	cmd.Flags().StringVar(&filterFlags.excludeNames, "exclude-names", "",
		"Comma-separated list of resource names to exclude")
}

// splitCSV splits a comma-separated string into a trimmed, non-empty slice.
// Returns nil if the input string is empty.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// buildFilterOptions converts the raw flag strings into a filter.Options struct.
func buildFilterOptions() filter.Options {
	return filter.Options{
		IncludeTypes: splitCSV(filterFlags.includeTypes),
		ExcludeTypes: splitCSV(filterFlags.excludeTypes),
		IncludeNames: splitCSV(filterFlags.includeNames),
		ExcludeNames: splitCSV(filterFlags.excludeNames),
	}
}
