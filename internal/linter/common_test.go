package linter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/linter"
)

// Create a file in a temporary folder, fill it with some content, and close it.
func createFileComplex(t *testing.T, file, content string) {
	t.Helper()

	require.NoError(t, os.WriteFile(file, []byte(content), 0o600))
}

// helper function to iterate the file until EOF, and return the lines read.
func ReadAll(t *testing.T, file *linter.Reader) []string {
	t.Helper()

	rows := make([]string, 0)

	require.NoError(t, file.Reset())

	for file.HasLines() {
		line, err := file.Next()
		rows = append(rows, line)

		require.NoError(t, err)
	}

	return rows
}

func CreateTempFile(t *testing.T, content ...string) string {
	t.Helper()

	file := TempFileName(t)

	require.NoError(t, os.WriteFile(file, []byte(strings.Join(content, "\n")), 0o600))

	return file
}

func TempFileName(t *testing.T) string {
	t.Helper()

	return filepath.Join(t.TempDir(), "test.txt")
}
