package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/owner/driftctl-diff/internal/diff"
)

type slackBlock struct {
	Type string      `json:"type"`
	Text *slackText  `json:"text,omitempty"`
	Fields []slackText `json:"fields,omitempty"`
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type slackPayload struct {
	Blocks []slackBlock `json:"blocks"`
}

func writeSlack(w io.Writer, result diff.Result) error {
	var blocks []slackBlock

	// Header
	headerText := ":white_check_mark: *No drift detected*"
	if result.HasDrift() {
		headerText = ":warning: *Terraform drift detected*"
	}
	blocks = append(blocks, slackBlock{
		Type: "section",
		Text: &slackText{Type: "mrkdwn", Text: headerText},
	})

	blocks = append(blocks, slackBlock{Type: "divider"})

	// Summary fields
	blocks = append(blocks, slackBlock{
		Type: "section",
		Fields: []slackText{
			{Type: "mrkdwn", Text: fmt.Sprintf("*Added:* %d", len(result.Added))},
			{Type: "mrkdwn", Text: fmt.Sprintf("*Removed:* %d", len(result.Removed))},
			{Type: "mrkdwn", Text: fmt.Sprintf("*Changed:* %d", len(result.Changed))},
		},
	})

	// Detail sections
	if len(result.Added) > 0 {
		keys := sortedKeys(result.Added)
		var buf bytes.Buffer
		buf.WriteString("*Added resources:*\n")
		for _, k := range keys {
			fmt.Fprintf(&buf, "• `%s`\n", k)
		}
		blocks = append(blocks, slackBlock{
			Type: "section",
			Text: &slackText{Type: "mrkdwn", Text: buf.String()},
		})
	}

	if len(result.Removed) > 0 {
		keys := sortedKeys(result.Removed)
		var buf bytes.Buffer
		buf.WriteString("*Removed resources:*\n")
		for _, k := range keys {
			fmt.Fprintf(&buf, "• `%s`\n", k)
		}
		blocks = append(blocks, slackBlock{
			Type: "section",
			Text: &slackText{Type: "mrkdwn", Text: buf.String()},
		})
	}

	if len(result.Changed) > 0 {
		keys := make([]string, 0, len(result.Changed))
		for k := range result.Changed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var buf bytes.Buffer
		buf.WriteString("*Changed resources:*\n")
		for _, k := range keys {
			fmt.Fprintf(&buf, "• `%s` (%d attribute(s))\n", k, len(result.Changed[k]))
		}
		blocks = append(blocks, slackBlock{
			Type: "section",
			Text: &slackText{Type: "mrkdwn", Text: buf.String()},
		})
	}

	payload := slackPayload{Blocks: blocks}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
