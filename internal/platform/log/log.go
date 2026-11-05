// Package log builds the process logger.
package log

import (
	"io"
	"log/slog"

	"github.com/georgifotev1/remr/internal/platform/config"
)

// New returns a logger writing to w in the configured format and level.
func New(w io.Writer, cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: cfg.Log.Level}

	var h slog.Handler = slog.NewJSONHandler(w, opts)
	if cfg.Log.Format == "text" {
		h = slog.NewTextHandler(w, opts)
	}
	return slog.New(h).With("service", cfg.Service, "version", cfg.Version, "env", cfg.Env)
}
