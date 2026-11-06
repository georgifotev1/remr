// Package db holds the Postgres pool and its settings.
package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/georgifotev1/remy/internal/platform/config"
)

// NewPool connects to Postgres and pings it, so a bad DB_URL fails at startup.
func NewPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	pc, err := poolConfig(cfg)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, pc)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

func poolConfig(cfg config.Config) (*pgxpool.Config, error) {
	pc, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		return nil, errors.New("invalid DB_URL") // pgx's error can echo the DSN, password included
	}
	pc.MaxConns = cfg.DB.MaxConns

	rp := pc.ConnConfig.RuntimeParams
	rp["application_name"] = cfg.Service
	rp["statement_timeout"] = millis(cfg.DB.StatementTimeout)
	rp["lock_timeout"] = millis(cfg.DB.LockTimeout)
	rp["idle_in_transaction_session_timeout"] = millis(cfg.DB.IdleInTxSessionTimeout)
	return pc, nil
}

func millis(d time.Duration) string {
	return strconv.FormatInt(d.Milliseconds(), 10)
}
