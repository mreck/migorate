package migorate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
)

type DB interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type Queries struct {
	FileData map[string]string
	DB       DB
}

func LoadQueriesFromDir(ctx context.Context, dir string, db DB) (*Queries, error) {
	files, err := FromDir(dir)
	if err != nil {
		return nil, fmt.Errorf("loading queries failed: %w", err)
	}

	return LoadQueries(ctx, files, db)
}

func LoadQueriesEmbedFS(ctx context.Context, fs embed.FS, db DB) (*Queries, error) {
	files, err := FromEmbedFS(fs)
	if err != nil {
		return nil, fmt.Errorf("loading queries failed: %w", err)
	}

	return LoadQueries(ctx, files, db)
}

func LoadQueries(ctx context.Context, files []SQLFile, db DB) (*Queries, error) {
	filedata := map[string]string{}

	for _, file := range files {
		stmt, err := db.PrepareContext(ctx, file.Content)
		if err != nil {
			return nil, fmt.Errorf("invalid query: %s: %w", file.Name, err)
		}
		stmt.Close()

		filedata[file.Name] = file.Content
	}

	return &Queries{filedata, db}, nil
}

func (q *Queries) Prepare(ctx context.Context, filename string) (*sql.Stmt, error) {
	query, ok := q.FileData[filename]
	if !ok {
		return nil, fmt.Errorf("query not found: %s", filename)
	}

	return q.DB.PrepareContext(ctx, query)
}

func (q *Queries) Exec(ctx context.Context, filename string, args ...any) (sql.Result, error) {
	stmt, err := q.Prepare(ctx, filename)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, args...)
}

func (q *Queries) Query(ctx context.Context, filename string, args ...any) (*sql.Rows, error) {
	stmt, err := q.Prepare(ctx, filename)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}

func (q *Queries) QueryRow(ctx context.Context, filename string, args ...any) (*sql.Row, error) {
	stmt, err := q.Prepare(ctx, filename)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, args...), nil
}
