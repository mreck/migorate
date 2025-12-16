package migrations

import (
	_ "embed"
)

var (
	//go:embed sqlite3.sql
	Sqlite3 string
)
