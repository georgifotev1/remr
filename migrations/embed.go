// Package migrations embeds the goose SQL migrations.
package migrations

import "embed"

// FS holds the embedded goose migrations.
//
//go:embed *.sql
var FS embed.FS
