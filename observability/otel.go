package observability

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// InitTracer configures OpenTelemetry tracing for the service and returns a shutdown function.
// It aims to maximize spans:
// - Sets global tracer provider with rich resource attributes (name, version, env).
// - Configures W3C TraceContext + Baggage propagators.
// - Wraps the default HTTP client transport to capture outbound HTTP spans.
// - Normalizes OTLP endpoint, supporting http/https schemes and insecure mode.
func InitTracer(ctx context.Context, serviceName string) func() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://otel-collector:4317" // default for local/docker
	}

	// Determine security mode by scheme and normalize endpoint to host:port for gRPC.
	useInsecure := true
	if strings.HasPrefix(endpoint, "https://") {
		useInsecure = false
	}
	endpoint = normalizeOTLPGRPCEndpoint(endpoint)

	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
	if useInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	// Resource attributes: service name, version, environment
	serviceVersion := getenvDefault("SERVICE_VERSION", "")
	serviceEnv := getenvDefault("OTEL_SERVICE_ENV", getenvDefault("SENTRY_ENV", ""))
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.DeploymentEnvironment(serviceEnv),
	)

	// Sampler: ratio from env (0..1). Defaults to 1.0 for maximum spans.
	ratio := parseRatio(getenvDefault("OTEL_SAMPLER_RATIO", "1"), 1.0)
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)

	// Global propagator: W3C TraceContext + Baggage
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	// Instrument default outbound HTTP client, idempotently
	if _, already := http.DefaultTransport.(*otelhttp.Transport); !already {
		http.DefaultTransport = otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithTracerProvider(tp),
			otelhttp.WithPropagators(otel.GetTextMapPropagator()),
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return r.Method + " " + r.URL.Host
			}),
		)
	}

	// Return shutdown function
	return func() {
		// Ensure provider shutdown flushes spans
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown tracer: %v", err)
		}
	}
}

var schemeRegex = regexp.MustCompile(`^(?i)(https?://)`)

// normalizeOTLPGRPCEndpoint removes http/https scheme prefixes if present.
// The gRPC exporter expects just host:port.
func normalizeOTLPGRPCEndpoint(s string) string {
	if s == "" {
		return s
	}
	return schemeRegex.ReplaceAllString(s, "")
}

// getenvDefault returns env var or fallback if empty.
func getenvDefault(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}

// parseRatio converts string to float64 in [0,1]; returns def on error/out of range.
func parseRatio(s string, def float64) float64 {
	if s == "" {
		return def
	}
	var r float64
	// simple parsing without extra deps
	for _, ch := range s {
		if !(ch == '.' || (ch >= '0' && ch <= '9')) {
			return def
		}
	}
	_, err := fmt.Sscan(s, &r)
	if err != nil || r < 0 || r > 1 {
		return def
	}
	return r
}
