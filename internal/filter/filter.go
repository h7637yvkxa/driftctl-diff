package filter

import (
	"strings"
)

// Options holds filtering criteria for drift results.
type Options struct {
	IncludeTypes []string
	ExcludeTypes []string
	IncludeNames []string
	ExcludeNames []string
}

// ResourceKey represents a parsed resource identifier.
type ResourceKey struct {
	Type string
	Name string
}

// ParseKey splits a "type.name" resource key into its components.
func ParseKey(key string) ResourceKey {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return ResourceKey{Type: key, Name: ""}
	}
	return ResourceKey{Type: parts[0], Name: parts[1]}
}

// Apply returns true if the resource key passes the filter options.
func Apply(key string, opts Options) bool {
	rk := ParseKey(key)

	if len(opts.ExcludeTypes) > 0 {
		for _, t := range opts.ExcludeTypes {
			if strings.EqualFold(rk.Type, t) {
				return false
			}
		}
	}

	if len(opts.IncludeTypes) > 0 {
		matched := false
		for _, t := range opts.IncludeTypes {
			if strings.EqualFold(rk.Type, t) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if len(opts.ExcludeNames) > 0 {
		for _, n := range opts.ExcludeNames {
			if strings.EqualFold(rk.Name, n) {
				return false
			}
		}
	}

	if len(opts.IncludeNames) > 0 {
		matched := false
		for _, n := range opts.IncludeNames {
			if strings.EqualFold(rk.Name, n) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
