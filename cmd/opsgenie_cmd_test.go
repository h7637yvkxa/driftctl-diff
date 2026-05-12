package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeOpsGenieState(t *testing.T, dir, name string, resources []map[string]interface{}) string {
	t.Helper()
	state := map[string]interface{}{
		"version": 4,
		"resources": resources,
	}
	data, _ := json.Marshal(state)
	p := filepath.Join(dir, name)
	_ = os.WriteFile(p, data, 0644)
	return p
}

func TestOpsGenieCmd_MissingAPIKey(t *testing.T) {
	os.Unsetenv("OPSGENIE_API_KEY")
	dir := t.TempDir()
	base := writeOpsGenieState(t, dir, "base.tfstate", nil)
	target := writeOpsGenieState(t, dir, "target.tfstate", nil)

	root := Execute()
	root.SetArgs([]string{"opsgenie", base, target})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when API key is missing")
	}
}

func TestOpsGenieCmd_NoDrift(t *testing.T) {
	var received bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = true
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	dir := t.TempDir()
	res := []map[string]interface{}{
		{"type": "aws_s3_bucket", "name": "my_bucket", "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
			"instances": []map[string]interface{}{{"attributes": map[string]interface{}{"id": "bucket-1"}}}},
	}
	base := writeOpsGenieState(t, dir, "base.tfstate", res)
	target := writeOpsGenieState(t, dir, "target.tfstate", res)

	root := Execute()
	root.SetArgs([]string{"opsgenie", "--api-key", "test-key", "--api-url", ts.URL, base, target})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !received {
		t.Error("expected OpsGenie endpoint to be called")
	}
}
