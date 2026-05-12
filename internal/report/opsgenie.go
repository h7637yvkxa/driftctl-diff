package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/owner/driftctl-diff/internal/diff"
)

type opsgeniePayload struct {
	Message     string            `json:"message"`
	Description string            `json:"description"`
	Priority    string            `json:"priority"`
	Tags        []string          `json:"tags"`
	Details     map[string]string `json:"details"`
	CreatedAt   string            `json:"createdAt"`
}

func writeOpsGenie(w io.Writer, result diff.Result, apiKey, apiURL string) error {
	if apiURL == "" {
		apiURL = "https://api.opsgenie.com/v2/alerts"
	}

	payload := buildOpsGeniePayload(result)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}

	_, err = fmt.Fprintf(w, "OpsGenie alert created (status %d)\n", resp.StatusCode)
	return err
}

func buildOpsGeniePayload(result diff.Result) opsgeniePayload {
	priority := "P5"
	if len(result.Added)+len(result.Removed)+len(result.Changed) > 0 {
		priority = "P3"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Added: %d, Removed: %d, Changed: %d\n",
		len(result.Added), len(result.Removed), len(result.Changed)))
	for k := range result.Changed {
		sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
	}

	tags := []string{"terraform", "drift"}
	if len(result.Added) > 0 {
		tags = append(tags, "added")
	}
	if len(result.Removed) > 0 {
		tags = append(tags, "removed")
	}

	return opsgeniePayload{
		Message:     fmt.Sprintf("Terraform drift detected: %d change(s)", len(result.Added)+len(result.Removed)+len(result.Changed)),
		Description: sb.String(),
		Priority:    priority,
		Tags:        tags,
		Details: map[string]string{
			"added":   fmt.Sprintf("%d", len(result.Added)),
			"removed": fmt.Sprintf("%d", len(result.Removed)),
			"changed": fmt.Sprintf("%d", len(result.Changed)),
		},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
