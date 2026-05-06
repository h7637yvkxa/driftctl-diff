package filter_test

import (
	"testing"

	"github.com/user/driftctl-diff/internal/filter"
)

func TestParseKey(t *testing.T) {
	rk := filter.ParseKey("aws_instance.web")
	if rk.Type != "aws_instance" || rk.Name != "web" {
		t.Fatalf("unexpected: %+v", rk)
	}

	rk2 := filter.ParseKey("malformed")
	if rk2.Type != "malformed" || rk2.Name != "" {
		t.Fatalf("unexpected: %+v", rk2)
	}
}

func TestApply_NoFilter(t *testing.T) {
	if !filter.Apply("aws_instance.web", filter.Options{}) {
		t.Fatal("expected key to pass empty filter")
	}
}

func TestApply_IncludeTypes(t *testing.T) {
	opts := filter.Options{IncludeTypes: []string{"aws_instance"}}
	if !filter.Apply("aws_instance.web", opts) {
		t.Fatal("expected aws_instance to be included")
	}
	if filter.Apply("aws_s3_bucket.data", opts) {
		t.Fatal("expected aws_s3_bucket to be excluded")
	}
}

func TestApply_ExcludeTypes(t *testing.T) {
	opts := filter.Options{ExcludeTypes: []string{"aws_s3_bucket"}}
	if filter.Apply("aws_s3_bucket.data", opts) {
		t.Fatal("expected aws_s3_bucket to be excluded")
	}
	if !filter.Apply("aws_instance.web", opts) {
		t.Fatal("expected aws_instance to pass")
	}
}

func TestApply_IncludeNames(t *testing.T) {
	opts := filter.Options{IncludeNames: []string{"web"}}
	if !filter.Apply("aws_instance.web", opts) {
		t.Fatal("expected 'web' to be included")
	}
	if filter.Apply("aws_instance.db", opts) {
		t.Fatal("expected 'db' to be excluded")
	}
}

func TestApply_ExcludeNames(t *testing.T) {
	opts := filter.Options{ExcludeNames: []string{"db"}}
	if filter.Apply("aws_instance.db", opts) {
		t.Fatal("expected 'db' to be excluded")
	}
	if !filter.Apply("aws_instance.web", opts) {
		t.Fatal("expected 'web' to pass")
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	opts := filter.Options{IncludeTypes: []string{"AWS_INSTANCE"}}
	if !filter.Apply("aws_instance.web", opts) {
		t.Fatal("expected case-insensitive match")
	}
}
