// Package app is the composition root: the only package that knows every dependency.
package app

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/georgifotev1/remr/internal/platform/config"
	"github.com/georgifotev1/remr/internal/platform/db"
)

// Deps holds the process-wide dependencies.
type Deps struct {
	Config config.Config
	Logger *slog.Logger
	Pool   *pgxpool.Pool
}

// Build creates the process dependencies.
func Build(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Deps, error) {
	pool, err := db.NewPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Deps{
		Config: cfg,
		Logger: logger,
		Pool:   pool,
	}, nil
}

// Close releases the process dependencies.
func (d *Deps) Close() {
	d.Pool.Close()
}
