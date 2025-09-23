package observability

import (
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
	"log/slog"
	"os"
)

func CreateLogger() *slog.Logger {
	return slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(os.Stdout, nil),
			slogsentry.Option{Level: slog.LevelError}.NewSentryHandler(),
		),
	)
}

func CreateLoggerWithTelemetry() *slog.Logger {
	return slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(os.Stdout, nil),
			otelslog.NewHandler("main", otelslog.WithLoggerProvider(global.GetLoggerProvider())),
			slogsentry.Option{Level: slog.LevelError}.NewSentryHandler(),
		),
	)
}
