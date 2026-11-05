package config

import "testing"

func TestLoadValid(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("DB_URL", "postgres://remr@localhost:5432/remr?sslmode=disable")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if c.Service != "remr" {
		t.Errorf("Service = %q, want remr", c.Service)
	}
	if c.HTTP.Addr != ":8080" {
		t.Errorf("HTTP.Addr = %q, want :8080", c.HTTP.Addr)
	}
	if c.DB.MaxConns != 10 {
		t.Errorf("DB.MaxConns = %d, want 10", c.DB.MaxConns)
	}
}

func TestLoadInvalidEnv(t *testing.T) {
	t.Setenv("APP_ENV", "nope")
	t.Setenv("DB_URL", "postgres://remr@localhost:5432/remr?sslmode=disable")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
}

func TestLoadMissingDBURL(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("DB_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want an error")
	}
}
