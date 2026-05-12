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

func TestWriteOpsGenie_NoDrift(t *testing.T) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		body, _ := io.ReadAll(r.Body)
		var p opsgeniePayload
		if err := json.Unmarshal(body, &p); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if p.Priority != "P5" {
			t.Errorf("expected P5 priority for no drift, got %s", p.Priority)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	result := diff.Result{}
	var out strings.Builder
	if err := writeOpsGenie(&out, result, "test-key", ts.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected HTTP call to OpsGenie")
	}
}

func TestWriteOpsGenie_WithDrift(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var p opsgeniePayload
		_ = json.Unmarshal(body, &p)
		if p.Priority != "P3" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	result := diff.Result{
		Added:   map[string]diff.Resource{"aws_s3_bucket.foo": {}},
		Removed: map[string]diff.Resource{},
		Changed: map[string]diff.Changed{},
	}
	var out strings.Builder
	if err := writeOpsGenie(&out, result, "key", ts.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteOpsGenie_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	var out strings.Builder
	err := writeOpsGenie(&out, diff.Result{}, "key", ts.URL)
	if err == nil {
		t.Fatal("expected error on server 500")
	}
}

func TestWriteOpsGenie_BadURL(t *testing.T) {
	var out strings.Builder
	err := writeOpsGenie(&out, diff.Result{}, "key", "http://127.0.0.1:0")
	if err == nil {
		t.Fatal("expected error on bad URL")
	}
}

func TestBuildOpsGeniePayload_Tags(t *testing.T) {
	result := diff.Result{
		Added:   map[string]diff.Resource{"r1": {}},
		Removed: map[string]diff.Resource{"r2": {}},
		Changed: map[string]diff.Changed{},
	}
	p := buildOpsGeniePayload(result)
	tagSet := map[string]bool{}
	for _, tag := range p.Tags {
		tagSet[tag] = true
	}
	if !tagSet["added"] {
		t.Error("expected 'added' tag")
	}
	if !tagSet["removed"] {
		t.Error("expected 'removed' tag")
	}
}
