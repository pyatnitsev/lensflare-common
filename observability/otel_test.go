package observability

import (
	"context"
	"net/http"
	"os"
	"testing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func TestInitTracer_Shutdown(t *testing.T) {
	// Use a localhost endpoint that shouldn't cause exporter creation to fail.
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	shutdown := InitTracer(context.Background(), "test-svc")
	if shutdown == nil {
		t.Fatalf("expected shutdown func")
	}
	// Should be safe to call
	shutdown()
}

func TestInitTracer_WrapsHTTPDefaultTransport_AndPropagator(t *testing.T) {
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
	shutdown := InitTracer(context.Background(), "svc")
	defer shutdown()

	if _, ok := http.DefaultTransport.(*otelhttp.Transport); !ok {
		t.Fatalf("expected http.DefaultTransport to be otelhttp.Transport")
	}

	p := otel.GetTextMapPropagator()
	if p == nil {
		t.Fatalf("expected global propagator to be set")
	}
	fields := map[string]bool{}
	for _, f := range p.Fields() {
		fields[f] = true
	}
	if !fields["traceparent"] {
		t.Fatalf("expected propagator to include traceparent field, got %#v", p.Fields())
	}
}

func TestInitTracer_EndpointNormalization_HTTP_HTTPS(t *testing.T) {
	for _, ep := range []string{"http://example:4317", "https://example:4317"} {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", ep)
		shutdown := InitTracer(context.Background(), "svc")
		if shutdown == nil {
			t.Fatalf("expected shutdown for %s", ep)
		}
		shutdown()
	}
}
