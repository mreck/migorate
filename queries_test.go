package migorate_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/mreck/migorate"
	"github.com/mreck/migorate/testutils"
	"github.com/mreck/migorate/testutils/sqlite3/migrations"
	"github.com/mreck/migorate/testutils/sqlite3/queries"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadQueriesFromEmbedFS(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	err = migorate.MigrateEmbedFS(ctx, "sqlite3", db, migrations.FS)
	require.NoError(t, err)

	files, err := migorate.FromEmbedFS(queries.FS)
	require.NoError(t, err)

	queries, err := migorate.LoadQueriesFromEmbedFS(ctx, files, db)
	assert.NoError(t, err)

	expected := map[string]string{
		"insert.sql":     testutils.RequireReadFile(t, "testutils/sqlite3/queries/insert.sql"),
		"select_all.sql": testutils.RequireReadFile(t, "testutils/sqlite3/queries/select_all.sql"),
		"select_one.sql": testutils.RequireReadFile(t, "testutils/sqlite3/queries/select_one.sql"),
	}
	assert.Equal(t, expected, queries.FileData)
}

func Test_Queries(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	err = migorate.MigrateEmbedFS(ctx, "sqlite3", db, migrations.FS)
	require.NoError(t, err)

	files, err := migorate.FromEmbedFS(queries.FS)
	require.NoError(t, err)

	queries, err := migorate.LoadQueriesFromEmbedFS(ctx, files, db)
	require.NoError(t, err)

	res, err := queries.Exec(ctx, "insert.sql", 1)
	assert.NoError(t, err)
	id, err := res.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	res, err = queries.Exec(ctx, "insert.sql", 2)
	assert.NoError(t, err)
	id, err = res.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), id)

	rows, err := queries.Query(ctx, "select_all.sql")
	assert.NoError(t, err)

	var ids []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		assert.NoError(t, err)
		ids = append(ids, id)
	}
	assert.Equal(t, []int64{1, 2}, ids)

	row, err := queries.QueryRow(ctx, "select_one.sql", 2)
	assert.NoError(t, err)
	assert.NoError(t, row.Err())

	var rowID int64
	err = row.Scan(&rowID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), rowID)
}
