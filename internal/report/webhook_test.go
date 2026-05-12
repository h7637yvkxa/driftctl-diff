package report

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
)

func newWebhookResult() diff.Result {
	return diff.Result{
		Added:   map[string]map[string]interface{}{"aws_s3_bucket.logs": {}},
		Removed: map[string]map[string]interface{}{},
		Changed: map[string]map[string][2]interface{}{
			"aws_instance.web": {"instance_type": {"t2.micro", "t3.small"}},
		},
	}
}

func TestWriteWebhook_Success(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	var buf strings.Builder
	err := writeWebhook(&buf, newWebhookResult(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var p webhookPayload
	if err := json.Unmarshal(received, &p); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if len(p.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(p.Added))
	}
	if len(p.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(p.Changed))
	}
	if p.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
	if !strings.Contains(buf.String(), ts.URL) {
		t.Errorf("output should mention URL, got: %s", buf.String())
	}
}

func TestWriteWebhook_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	var buf strings.Builder
	err := writeWebhook(&buf, newWebhookResult(), ts.URL)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestWriteWebhook_BadURL(t *testing.T) {
	var buf strings.Builder
	err := writeWebhook(&buf, newWebhookResult(), "http://127.0.0.1:0/no-server")
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestBuildWebhookPayload_NoDrift(t *testing.T) {
	result := diff.Result{
		Added:   map[string]map[string]interface{}{},
		Removed: map[string]map[string]interface{}{},
		Changed: map[string]map[string][2]interface{}{},
	}
	p := buildWebhookPayload(result)
	if len(p.Added) != 0 || len(p.Removed) != 0 || len(p.Changed) != 0 {
		t.Error("expected empty slices for no-drift result")
	}
}
