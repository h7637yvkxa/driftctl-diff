package ignore

import (
	"bufio"
	"os"
	"strings"
)

// Rule represents a single ignore rule.
type Rule struct {
	ResourceType string
	ResourceName string
	Attribute    string
}

// Set holds a collection of ignore rules.
type Set struct {
	rules []Rule
}

// ParseFile reads an ignore file and returns a Set of rules.
// Each non-blank, non-comment line should be in the form:
//
//	<type>.<name>[.<attribute>]
func ParseFile(path string) (*Set, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	set := &Set{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rule, err := parseLine(line)
		if err != nil {
			continue // skip malformed lines
		}
		set.rules = append(set.rules, rule)
	}
	return set, scanner.Err()
}

func parseLine(line string) (Rule, error) {
	parts := strings.SplitN(line, ".", 3)
	if len(parts) < 2 {
		return Rule{}, fmt.Errorf("invalid rule: %s", line)
	}
	r := Rule{
		ResourceType: parts[0],
		ResourceName: parts[1],
	}
	if len(parts) == 3 {
		r.Attribute = parts[2]
	}
	return r, nil
}

// Matches returns true if the given type, name, and attribute match any rule.
// Pass an empty attribute to check for a whole-resource ignore.
func (s *Set) Matches(resourceType, resourceName, attribute string) bool {
	for _, r := range s.rules {
		typeMatch := r.ResourceType == "*" || r.ResourceType == resourceType
		nameMatch := r.ResourceName == "*" || r.ResourceName == resourceName
		attrMatch := r.Attribute == "" || r.Attribute == "*" || r.Attribute == attribute
		if typeMatch && nameMatch && attrMatch {
			return true
		}
	}
	return false
}
