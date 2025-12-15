package migorate_test

import (
	"testing"

	"github.com/mreck/migorate"
	"github.com/mreck/migorate/testutils"
	"github.com/mreck/migorate/testutils/sqlite3/migrations"

	"github.com/stretchr/testify/assert"
)

func Test_FromEmbedFS(t *testing.T) {
	m, err := migorate.FromEmbedFS(migrations.FS)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	assert.Equal(t, m[0].Name, "001.sql")
	assert.Equal(t, m[1].Name, "002.sql")

	assert.Equal(t, m[0].Content, testutils.RequireReadFile(t, "testutils/sqlite3/migrations/001.sql"))
	assert.Equal(t, m[1].Content, testutils.RequireReadFile(t, "testutils/sqlite3/migrations/002.sql"))
}

func Test_FromDir(t *testing.T) {
	m, err := migorate.FromDir("testutils/sqlite3/migrations")
	assert.NoError(t, err)
	assert.NotNil(t, m)

	assert.Equal(t, m[0].Name, "001.sql")
	assert.Equal(t, m[1].Name, "002.sql")

	assert.Equal(t, m[0].Content, testutils.RequireReadFile(t, "testutils/sqlite3/migrations/001.sql"))
	assert.Equal(t, m[1].Content, testutils.RequireReadFile(t, "testutils/sqlite3/migrations/002.sql"))
}
