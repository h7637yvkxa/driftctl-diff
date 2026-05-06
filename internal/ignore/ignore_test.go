package ignore

import (
	"os"
	"testing"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.driftignore")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Valid(t *testing.T) {
	path := writeTempIgnore(t, "# comment\naws_s3_bucket.my_bucket\naws_instance.web.tags\n")
	set, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(set.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(set.rules))
	}
}

func TestParseFile_Missing(t *testing.T) {
	_, err := ParseFile("/nonexistent/.driftignore")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMatches_WholeResource(t *testing.T) {
	path := writeTempIgnore(t, "aws_s3_bucket.my_bucket\n")
	set, _ := ParseFile(path)
	if !set.Matches("aws_s3_bucket", "my_bucket", "") {
		t.Error("expected match for whole-resource rule")
	}
	if set.Matches("aws_s3_bucket", "other_bucket", "") {
		t.Error("unexpected match for different name")
	}
}

func TestMatches_Attribute(t *testing.T) {
	path := writeTempIgnore(t, "aws_instance.web.tags\n")
	set, _ := ParseFile(path)
	if !set.Matches("aws_instance", "web", "tags") {
		t.Error("expected attribute match")
	}
	if set.Matches("aws_instance", "web", "ami") {
		t.Error("unexpected match for different attribute")
	}
}

func TestMatches_Wildcard(t *testing.T) {
	path := writeTempIgnore(t, "aws_s3_bucket.*.tags\n")
	set, _ := ParseFile(path)
	if !set.Matches("aws_s3_bucket", "any_bucket", "tags") {
		t.Error("expected wildcard name match")
	}
}
