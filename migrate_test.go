package migorate_test

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/mreck/migorate"
	"github.com/mreck/migorate/testutils/sqlite3/migrations"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_MigrateFS(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	assert.Nil(t, err)

	m, err := migorate.FromEmbedFS(migrations.FS)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for i := range 3 {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			err := migorate.Migrate(ctx, "sqlite3", db, m)
			assert.Nil(t, err)

			tables, err := getTablesSqlite(db)
			assert.Nil(t, err)
			assert.Equal(t, []string{"migrations", "test_1", "test_2"}, tables)
		})
	}

	m = append(m, migorate.SQLFile{"3", "CREATE TABLE test_3 (id INT)"})

	for i := range 3 {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			err := migorate.Migrate(ctx, "sqlite3", db, m)
			assert.Nil(t, err)

			tables, err := getTablesSqlite(db)
			assert.Nil(t, err)
			assert.Equal(t, []string{"migrations", "test_1", "test_2", "test_3"}, tables)
		})
	}
}

func getTablesSqlite(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`
		SELECT name
		FROM sqlite_schema
		WHERE type = 'table'
		AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		tables = append(tables, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	sort.Strings(tables)

	return tables, nil
}
