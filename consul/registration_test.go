package consul

import (
	"strings"
	"testing"
)

func TestRegisterService_MissingConsulAddr(t *testing.T) {
	// Unset all relevant envs
	t.Setenv("CONSUL_HTTP_ADDR", "")
	t.Setenv("HTTP_PORT", "")
	t.Setenv("CONSUL_SERVICE_ID", "")
	t.Setenv("CONSUL_SERVICE_NAME", "")
	t.Setenv("HOST_ADDRESS", "")

	err := RegisterService()
	if err == nil {
		t.Fatalf("expected error for missing CONSUL_HTTP_ADDR, got nil")
	}
	if got := err.Error(); got != "CONSUL_HTTP_ADDR not set" {
		t.Fatalf("unexpected error: %q", got)
	}
}

func TestRegisterService_MissingHTTPPort(t *testing.T) {
	t.Setenv("CONSUL_HTTP_ADDR", "http://127.0.0.1:8500")
	t.Setenv("HTTP_PORT", "")
	t.Setenv("CONSUL_SERVICE_ID", "")
	t.Setenv("CONSUL_SERVICE_NAME", "")
	t.Setenv("HOST_ADDRESS", "")

	err := RegisterService()
	if err == nil {
		t.Fatalf("expected error for missing HTTP_PORT, got nil")
	}
	if got := err.Error(); got != "HTTP_PORT not set" {
		t.Fatalf("unexpected error: %q", got)
	}
}

func TestRegisterService_InvalidHTTPPort(t *testing.T) {
	t.Setenv("CONSUL_HTTP_ADDR", "http://127.0.0.1:8500")
	t.Setenv("HTTP_PORT", "abc")
	t.Setenv("CONSUL_SERVICE_ID", "")
	t.Setenv("CONSUL_SERVICE_NAME", "")
	t.Setenv("HOST_ADDRESS", "")

	err := RegisterService()
	if err == nil {
		t.Fatalf("expected error for invalid HTTP_PORT, got nil")
	}
	if got := err.Error(); !strings.HasPrefix(got, "Invalid HTTP_PORT") {
		t.Fatalf("unexpected error: %q", got)
	}
}

func TestRegisterService_MissingServiceID(t *testing.T) {
	t.Setenv("CONSUL_HTTP_ADDR", "http://127.0.0.1:8500")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("CONSUL_SERVICE_ID", "")
	t.Setenv("CONSUL_SERVICE_NAME", "")
	t.Setenv("HOST_ADDRESS", "")

	err := RegisterService()
	if err == nil {
		t.Fatalf("expected error for missing CONSUL_SERVICE_ID, got nil")
	}
	if got := err.Error(); got != "CONSUL_SERVICE_ID not set" {
		t.Fatalf("unexpected error: %q", got)
	}
}

func TestRegisterService_MissingServiceName(t *testing.T) {
	t.Setenv("CONSUL_HTTP_ADDR", "http://127.0.0.1:8500")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("CONSUL_SERVICE_ID", "svc-1")
	t.Setenv("CONSUL_SERVICE_NAME", "")
	t.Setenv("HOST_ADDRESS", "")

	err := RegisterService()
	if err == nil {
		t.Fatalf("expected error for missing CONSUL_SERVICE_NAME, got nil")
	}
	if got := err.Error(); got != "CONSUL_SERVICE_NAME not set" {
		t.Fatalf("unexpected error: %q", got)
	}
}

func TestRegisterService_MissingHostAddress(t *testing.T) {
	t.Setenv("CONSUL_HTTP_ADDR", "http://127.0.0.1:8500")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("CONSUL_SERVICE_ID", "svc-1")
	t.Setenv("CONSUL_SERVICE_NAME", "svc-name")
	t.Setenv("HOST_ADDRESS", "")

	err := RegisterService()
	if err == nil {
		t.Fatalf("expected error for missing HOST_ADDRESS, got nil")
	}
	if got := err.Error(); got != "HOST_ADDRESS not set" {
		t.Fatalf("unexpected error: %q", got)
	}
}
