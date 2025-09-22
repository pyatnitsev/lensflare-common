package db

import (
	"testing"
)

func TestInit_MissingPGURL(t *testing.T) {
	// Ensure PG_URL is unset for this test case
	t.Setenv("PG_URL", "")

	db, err := Init()
	if err == nil {
		t.Fatalf("expected error when PG_URL is missing, got nil")
	}
	if db != nil {
		t.Fatalf("expected returned *gorm.DB to be nil when error occurs")
	}
	if got, want := err.Error(), "environment variable PG_URL is not set"; got != want {
		// Match exact message to ensure clear diagnostics remain stable
		t.Fatalf("unexpected error: got %q, want %q", got, want)
	}
}
