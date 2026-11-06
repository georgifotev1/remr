// Package config loads the typed process configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config is parsed once at startup and passed to constructors by value.
type Config struct {
	Service string `env:"SERVICE_NAME" envDefault:"remy"`
	Env     string `env:"APP_ENV,required"`
	Version string

	HTTP HTTP `envPrefix:"HTTP_"`
	DB   DB   `envPrefix:"DB_"`
	Log  Log  `envPrefix:"LOG_"`
}

// HTTP holds the API server settings.
type HTTP struct {
	Addr            string        `env:"ADDR" envDefault:":8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"30s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"15s"`
}

// DB holds the Postgres pool settings. URL contains the password: never log it.
type DB struct {
	URL                    string        `env:"URL,notEmpty"`
	MaxConns               int32         `env:"MAX_CONNS" envDefault:"10"`
	StatementTimeout       time.Duration `env:"STATEMENT_TIMEOUT" envDefault:"5s"`
	LockTimeout            time.Duration `env:"LOCK_TIMEOUT" envDefault:"2s"`
	IdleInTxSessionTimeout time.Duration `env:"IDLE_IN_TX_SESSION_TIMEOUT" envDefault:"10s"`
}

// Log holds the logger settings.
type Log struct {
	Format string     `env:"FORMAT" envDefault:"json"`
	Level  slog.Level `env:"LEVEL" envDefault:"info"`
}

// Load parses the environment and validates the result.
func Load() (Config, error) {
	c, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return c, c.Validate()
}

// Validate reports every invalid setting at once.
func (c Config) Validate() error {
	var errs []error
	switch c.Env {
	case "development", "staging", "production":
	default:
		errs = append(errs, fmt.Errorf("APP_ENV must be development, staging or production, got %q", c.Env))
	}
	switch c.Log.Format {
	case "json", "text":
	default:
		errs = append(errs, fmt.Errorf("LOG_FORMAT must be json or text, got %q", c.Log.Format))
	}
	if c.HTTP.ReadTimeout <= 0 {
		errs = append(errs, errors.New("HTTP_READ_TIMEOUT must be > 0"))
	}
	if c.HTTP.WriteTimeout <= 0 {
		errs = append(errs, errors.New("HTTP_WRITE_TIMEOUT must be > 0"))
	}
	if c.HTTP.IdleTimeout <= 0 {
		errs = append(errs, errors.New("HTTP_IDLE_TIMEOUT must be > 0"))
	}
	if c.HTTP.ShutdownTimeout <= 0 {
		errs = append(errs, errors.New("HTTP_SHUTDOWN_TIMEOUT must be > 0"))
	}
	if c.DB.MaxConns <= 0 {
		errs = append(errs, errors.New("DB_MAX_CONNS must be > 0"))
	}
	if c.DB.StatementTimeout <= 0 {
		errs = append(errs, errors.New("DB_STATEMENT_TIMEOUT must be > 0"))
	}
	if c.DB.LockTimeout <= 0 {
		errs = append(errs, errors.New("DB_LOCK_TIMEOUT must be > 0"))
	}
	if c.DB.IdleInTxSessionTimeout <= 0 {
		errs = append(errs, errors.New("DB_IDLE_IN_TX_SESSION_TIMEOUT must be > 0"))
	}
	return errors.Join(errs...)
}
