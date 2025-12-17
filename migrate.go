package migorate

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/mreck/migorate/migrations"
)

type Migration struct {
	Key    string
	Script string
}

var (
	ErrDBNotSupported = errors.New("database not supported")
)

func MigrateFromDir(ctx context.Context, driverName string, db *sql.DB, dir string) error {
	files, err := FromDir(dir)
	if err != nil {
		return fmt.Errorf("creating migrations failed: %w", err)
	}

	return Migrate(ctx, driverName, db, files)
}

func MigrateEmbedFS(ctx context.Context, driverName string, db *sql.DB, fs embed.FS) error {
	files, err := FromEmbedFS(fs)
	if err != nil {
		return fmt.Errorf("creating migrations failed: %w", err)
	}

	return Migrate(ctx, driverName, db, files)
}

func Migrate(ctx context.Context, driverName string, db *sql.DB, files []SQLFile) error {
	var initQuery string

	switch driverName {
	case "sqlite3":
		initQuery = migrations.Sqlite3
	default:
		return ErrDBNotSupported
	}

	err := db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("pinging db failed: %w", err)
	}

	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	tx, err := db.BeginTx(cCtx, nil)
	if err != nil {
		return fmt.Errorf("creating transaction failed: %w", err)
	}

	_, err = tx.Exec(initQuery)
	if err != nil {
		return fmt.Errorf("running migration table script failed: %w", err)
	}

	for _, file := range files {
		var n int
		err := tx.QueryRow(`SELECT COUNT(*) FROM "migrations" WHERE "filename" = ?`, file.Name).Scan(&n)
		if err != nil {
			return fmt.Errorf("selecting from migrations table failed: %w", err)
		}
		if n > 0 {
			continue
		}

		_, err = tx.Exec(file.Content)
		if err != nil {
			return fmt.Errorf("running migration failed: %s: %w", file.Name, err)
		}

		_, err = tx.Exec(`INSERT INTO "migrations" ("filename") VALUES (?)`, file.Name)
		if err != nil {
			return fmt.Errorf("updating migrations table failed: %s: %w", file.Name, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("comitting migration transaction failed: %w", err)
	}

	return nil
}
