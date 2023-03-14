package file_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/file"
)

// Create a file in a temporary folder, fill it with some content, and close it.
func createFile(t *testing.T, file, content string) {
	t.Helper()

	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestFileHandler(t *testing.T) {
	t.Parallel()

	// Create a file with some contents
	filePath := filepath.Join(t.TempDir(), "test.txt")

	content := []string{
		"Line 1",
		"Line 2",
		"Line 3",
	}

	createFile(t, filePath, strings.Join(content, "\n"))

	// Create a file manager
	handler := file.Manager{Name: filePath}

	require.NoError(t, handler.Open())

	// Read all lines into a slice
	rows := make([]string, 0)

	for i := 0; handler.HasLines(); i++ {
		line, err := handler.Next()
		rows = append(rows, line)

		require.NoError(t, err)
	}

	// Check that the contents are the same
	require.Equal(t, content, rows)
}

func TestFileHandler_Error(t *testing.T) {
	t.Parallel()

	// Create a file with some contents
	filePath := filepath.Join(t.TempDir(), "test.txt")

	// Create a file manager
	handler := file.Manager{Name: filePath}

	// Require an error when opening a non-existing file
	require.Error(t, handler.Open())
	// Require an error when closing a non-opened file
	require.Error(t, handler.Close())

	// Create the file
	createFile(t, filePath, "no content")

	// Require no error when opening an existing file
	require.NoError(t, handler.Open())

	require.NoError(t, handler.Close())

	// Try to read the file
	_, err := handler.Next()
	require.Error(t, err)
}