package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type recordingTransport struct{ calls atomic.Int64 }

func (t *recordingTransport) Configure(options sentry.ClientOptions) {}
func (t *recordingTransport) SendEvent(event *sentry.Event) {
	t.calls.Add(1)
}
func (t *recordingTransport) Flush(timeout time.Duration) bool          { return true }
func (t *recordingTransport) FlushWithContext(ctx context.Context) bool { return true }
func (t *recordingTransport) Close()                                    {}

func TestErrorMiddleware_RespondsWithErrorJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorMiddleware())
	r.GET("/boom", func(c *gin.Context) {
		_ = c.Error(errors.New("boom"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/boom", nil)
	r.ServeHTTP(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON body, got %q: %v", rec.Body.String(), err)
	}
	if got, ok := body["error"].(string); !ok || got == "" {
		t.Fatalf("expected error field in JSON, got: %#v", body)
	}
}

func TestErrorMiddleware_CapturesToSentry_WhenHubInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a Sentry client with recording transport to avoid network
	rt := &recordingTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1", // any parseable DSN
		Transport: rt,
	})
	if err != nil {
		t.Fatalf("sentry.NewClient error: %v", err)
	}
	hub := sentry.NewHub(client, sentry.NewScope())

	r := gin.New()
	r.Use(ErrorMiddleware())
	r.GET("/boom", func(c *gin.Context) {
		_ = c.Error(errors.New("boom"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/boom", nil)
	// Attach sentry hub to the request context so middleware can pick it up
	ctx := sentry.SetHubOnContext(req.Context(), hub)
	req = req.WithContext(ctx)

	r.ServeHTTP(rec, req)

	if rt.calls.Load() == 0 {
		t.Fatalf("expected middleware to capture exception to Sentry, but transport was not called")
	}
}
