package report

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
)

// cycloneDX represents a minimal CycloneDX-style BOM repurposed for drift reporting.
type cycloneDXReport struct {
	XMLName    xml.Name        `xml:"bom"`
	XMLNS      string          `xml:"xmlns,attr"`
	Version    int             `xml:"version,attr"`
	Metadata   cdxMetadata     `xml:"metadata"`
	Components []cdxComponent  `xml:"components>component"`
}

type cdxMetadata struct {
	Timestamp string     `xml:"timestamp"`
	Tools     []cdxTool  `xml:"tools>tool"`
}

type cdxTool struct {
	Vendor  string `xml:"vendor"`
	Name    string `xml:"name"`
	Version string `xml:"version"`
}

type cdxComponent struct {
	Type        string          `xml:"type,attr"`
	Name        string          `xml:"name"`
	Description string          `xml:"description,omitempty"`
	Properties  []cdxProperty   `xml:"properties>property"`
}

type cdxProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func writeCycloneDX(w io.Writer, result diff.Result) error {
	var components []cdxComponent

	for key := range result.Added {
		components = append(components, cdxComponent{
			Type:        "library",
			Name:        key,
			Description: "added resource (present in source, missing in target)",
			Properties:  []cdxProperty{{Name: "drift.status", Value: "added"}},
		})
	}

	for key := range result.Removed {
		components = append(components, cdxComponent{
			Type:        "library",
			Name:        key,
			Description: "removed resource (missing in source, present in target)",
			Properties:  []cdxProperty{{Name: "drift.status", Value: "removed"}},
		})
	}

	for key, changes := range result.Changed {
		props := []cdxProperty{{Name: "drift.status", Value: "changed"}}
		for _, ch := range changes {
			props = append(props, cdxProperty{
				Name:  fmt.Sprintf("drift.attr.%s", ch.Attribute),
				Value: fmt.Sprintf("%v -> %v", ch.OldValue, ch.NewValue),
			})
		}
		components = append(components, cdxComponent{
			Type:        "library",
			Name:        key,
			Description: "changed resource",
			Properties:  props,
		})
	}

	report := cycloneDXReport{
		XMLNS:   "http://cyclonedx.org/schema/bom/1.4",
		Version: 1,
		Metadata: cdxMetadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Tools:     []cdxTool{{Vendor: "driftctl-diff", Name: "driftctl-diff", Version: "dev"}},
		},
		Components: components,
	}

	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("cyclonedx encode: %w", err)
	}
	return enc.Flush()
}
