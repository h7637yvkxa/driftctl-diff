package state_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftctl-diff/internal/state"
)

func writeTempState(t *testing.T, s state.TerraformState) string {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal state: %v", err)
	}
	tmp := filepath.Join(t.TempDir(), "terraform.tfstate")
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		t.Fatalf("failed to write temp state: %v", err)
	}
	return tmp
}

func TestParseStateFile_Valid(t *testing.T) {
	input := state.TerraformState{
		Version:   4,
		TFVersion: "1.5.0",
		Resources: []state.Resource{
			{Type: "aws_instance", Name: "web", Provider: "aws", Attributes: map[string]interface{}{"ami": "ami-123"}},
		},
	}
	path := writeTempState(t, input)

	got, err := state.ParseStateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Version != 4 {
		t.Errorf("expected version 4, got %d", got.Version)
	}
	if len(got.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(got.Resources))
	}
}

func TestParseStateFile_Missing(t *testing.T) {
	_, err := state.ParseStateFile("/nonexistent/path.tfstate")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestIndexResources(t *testing.T) {
	s := &state.TerraformState{
		Resources: []state.Resource{
			{Type: "aws_s3_bucket", Name: "assets"},
			{Type: "aws_instance", Name: "app"},
		},
	}
	index := state.IndexResources(s)
	if _, ok := index["aws_s3_bucket.assets"]; !ok {
		t.Error("expected key aws_s3_bucket.assets in index")
	}
	if _, ok := index["aws_instance.app"]; !ok {
		t.Error("expected key aws_instance.app in index")
	}
	if len(index) != 2 {
		t.Errorf("expected 2 entries, got %d", len(index))
	}
}
