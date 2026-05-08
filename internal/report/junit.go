package report

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/owner/driftctl-diff/internal/diff"
)

type junitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []junitTestSuite `xml:"testsuite"`
}

type junitTestSuite struct {
	Name     string          `xml:"name,attr"`
	Tests    int             `xml:"tests,attr"`
	Failures int             `xml:"failures,attr"`
	Cases    []junitTestCase `xml:"testcase"`
}

type junitTestCase struct {
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Failure   *junitFailure `xml:"failure,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

func writeJUnit(w io.Writer, result diff.Result) error {
	var cases []junitTestCase

	for key, res := range result.Added {
		cases = append(cases, junitTestCase{
			Name:      fmt.Sprintf("%s exists in target", key),
			Classname: res.Type,
			Failure: &junitFailure{
				Message: "resource added in target",
				Text:    fmt.Sprintf("Resource %q is present in target but not in source.", key),
			},
		})
	}

	for key, res := range result.Removed {
		cases = append(cases, junitTestCase{
			Name:      fmt.Sprintf("%s missing from target", key),
			Classname: res.Type,
			Failure: &junitFailure{
				Message: "resource removed from target",
				Text:    fmt.Sprintf("Resource %q is present in source but missing from target.", key),
			},
		})
	}

	for key, ch := range result.Changed {
		cases = append(cases, junitTestCase{
			Name:      fmt.Sprintf("%s has attribute drift", key),
			Classname: ch.Source.Type,
			Failure: &junitFailure{
				Message: fmt.Sprintf("%d attribute(s) changed", len(ch.Attributes)),
				Text:    fmt.Sprintf("Resource %q has %d changed attribute(s).", key, len(ch.Attributes)),
			},
		})
	}

	failures := len(result.Added) + len(result.Removed) + len(result.Changed)
	suite := junitTestSuite{
		Name:     "driftctl-diff",
		Tests:    len(cases),
		Failures: failures,
		Cases:    cases,
	}

	suites := junitTestSuites{Suites: []junitTestSuite{suite}}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(suites); err != nil {
		return fmt.Errorf("junit encode: %w", err)
	}
	return enc.Flush()
}
