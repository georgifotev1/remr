package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/georgifotev1/remr/internal/platform/config"
	"github.com/georgifotev1/remr/internal/platform/log"
	"github.com/georgifotev1/remr/internal/platform/spa"
	webassets "github.com/georgifotev1/remr/web"
)

type pinger interface {
	Ping(ctx context.Context) error
}

// RunAPI starts the HTTP API and serves the embedded SPA until ctx is cancelled.
func RunAPI(ctx context.Context, version string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cfg.Version = version

	logger := log.New(os.Stderr, cfg)
	slog.SetDefault(logger)

	deps, err := Build(ctx, cfg, logger)
	if err != nil {
		return err
	}
	defer deps.Close()

	server := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           newHandler(deps),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()
	logger.Info("listening", "addr", cfg.HTTP.Addr)

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("api server: %w", err)
		}
	case <-ctx.Done():
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		if err := <-errCh; !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("api server: %w", err)
		}
	}

	logger.Info("stopped cleanly")
	return nil
}

func newHandler(d *Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", healthHandler(d.Config.Version))
	mux.HandleFunc("GET /api/v1/ready", readyHandler(d.Pool))
	mux.Handle("/", spa.Handler(webassets.Dist()))
	return mux
}

func healthHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"version": version,
		})
	}
}

func readyHandler(p pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := p.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
