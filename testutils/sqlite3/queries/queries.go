package queries

import (
	"embed"
)

var (
	//go:embed *.sql
	FS embed.FS
)
