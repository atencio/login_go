// Package migrations embebe los scripts SQL que se ejecutan al iniciar la API.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
