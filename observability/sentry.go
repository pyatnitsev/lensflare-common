package observability

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
)

// InitSentry инициализирует Sentry и возвращает функцию graceful shutdown
func InitSentry() func() {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		log.Println("[Sentry] Not configured")
		return func() {}
	}
	sampleRate := 1.0
	if sr := os.Getenv("SENTRY_SAMPLE_RATE"); sr != "" {
		if val, err := strconv.ParseFloat(sr, 64); err == nil {
			// Clamp to [0,1]
			if !math.IsNaN(val) && !math.IsInf(val, 0) {
				if val < 0 {
					val = 0
				} else if val > 1 {
					val = 1
				}
				sampleRate = val
			}
		}
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: sampleRate,
		Environment:      os.Getenv("SENTRY_ENV"),
		Release:          os.Getenv("SENTRY_RELEASE"),
	}); err != nil {
		// Do not crash host application from a library
		log.Printf("[Sentry] init failed: %v", err)
		return func() {}
	}
	log.Println("[Sentry] Initialization complete")
	return func() {
		// Allow flush timeout override via env, fallback to 2s
		flushMs := 2000
		if v := os.Getenv("SENTRY_FLUSH_MS"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				flushMs = n
			}
		}
		sentry.Flush(time.Duration(flushMs) * time.Millisecond)
	}
}
