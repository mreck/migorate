package testutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func RequireReadFile(t *testing.T, filename string) string {
	b, err := os.ReadFile(filename)
	require.NoError(t, err)
	return string(b)
}
