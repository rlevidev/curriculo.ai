package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	// stdout JSON para coletor do host/Docker, sem agregador dedicado por enquanto
	baseLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.UTC().Format(time.RFC3339))
				}
			}
			return a
		},
	})).With("service", "curriculo-api")
)

// main reads PORT (defaults to 8080), starts the server, and handles ListenAndServe fatal errors.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	baseLogger.Info("starting", "op", "startup", "port", port)
	server := setupServer(port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		baseLogger.Error("server failed", "op", "startup", "error", err.Error())
		os.Exit(1)
	}
}
