package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
)

type webhookPayload struct {
	Timestamp  string            `json:"timestamp"`
	Summary    string            `json:"summary"`
	Added      []string          `json:"added,omitempty"`
	Removed    []string          `json:"removed,omitempty"`
	Changed    []string          `json:"changed,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

func writeWebhook(w io.Writer, result diff.Result, url string) error {
	payload := buildWebhookPayload(result)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body)) //nolint:gosec
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}

	_, err = fmt.Fprintf(w, "webhook delivered to %s (status %d)\n", url, resp.StatusCode)
	return err
}

func buildWebhookPayload(result diff.Result) webhookPayload {
	p := webhookPayload{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Summary:    result.Summary(),
		Attributes: make(map[string]string),
	}
	for k := range result.Added {
		p.Added = append(p.Added, k)
	}
	for k := range result.Removed {
		p.Removed = append(p.Removed, k)
	}
	for k, attrs := range result.Changed {
		p.Changed = append(p.Changed, k)
		for attr := range attrs {
			p.Attributes[k+"."+attr] = ""
		}
	}
	return p
}
