package observability

import (
	"testing"
)

// Test that InitSentry returns a no-op shutdown when SENTRY_DSN is empty
func TestInitSentry_NoDSN(t *testing.T) {
	// Ensure DSN is not set so init is skipped
	t.Setenv("SENTRY_DSN", "")
	// Even if sample rate/env/release are set, without DSN it should skip
	t.Setenv("SENTRY_SAMPLE_RATE", "not-a-float")
	t.Setenv("SENTRY_ENV", "test")
	t.Setenv("SENTRY_RELEASE", "1.2.3")

	shutdown := InitSentry()
	if shutdown == nil {
		t.Fatalf("expected non-nil shutdown function")
	}
	// Calling shutdown should not panic even if Sentry wasn't initialized
	// (function should be a no-op)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("shutdown panicked: %v", r)
		}
	}()
	shutdown()
}
