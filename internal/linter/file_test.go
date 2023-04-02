package linter_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/linter"
)

// Test the File by
// 1. Creating a file with some contents
// 2. Creating a File to interact with the file
// 3. Reading the file line by line
// 4. Checking that the contents are the same.
func TestFile(t *testing.T) {
	t.Parallel()

	// Create a file with some contents
	filePath := filepath.Join(t.TempDir(), "test.txt")

	content := []string{
		"Line 1",
		"Line 2",
		"Line 3",
	}

	createFile(t, filePath, strings.Join(content, "\n"))

	// Create a file File
	file, err := linter.NewFile(filePath)

	// No error should occur when opening
	require.NoError(t, err)

	// Read all lines into a slice
	rows := iterateFile(t, file)

	// Check that the contents are the same
	require.Equal(t, content, rows)
}

func TestFile_Write(t *testing.T) {
	t.Parallel()

	// Create a file with no content
	filePath := filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "")

	// Create a file File
	file, err := linter.NewFile(filePath)

	// No error should occur when opening
	require.NoError(t, err)

	content := []string{
		"Line 1",
		"Line 2",
		"Line 3",
	}

	require.NoError(t, file.Write(content...))

	// Read all lines into a slice
	rows := iterateFile(t, file)

	// Check that the contents are the same
	// The last line is empty because the file ends with a newline
	require.Equal(t, content, rows[:len(rows)-1])
}

// Test the File by
// 1. Creating an empty folder
// 2. Creating a File to interact with a non-existing file
// 3. Checking that an error occurs when opening the file.
func TestFile_ErrorFileNotExisting(t *testing.T) {
	t.Parallel()

	// Non-existing file
	filePath := filepath.Join(t.TempDir(), "test.txt")

	// Create a File pointing to a non-existing file
	_, err := linter.NewFile(filePath)

	// Require an error when opening a non-existing file
	require.Error(t, err)
}

func TestFile_ErrorCloseUnopenedFile(t *testing.T) {
	t.Parallel()

	// Non-existing file
	filePath := filepath.Join(t.TempDir(), "test.txt")

	// Create a File pointing to a non-existing file
	file, _ := linter.NewFile(filePath)

	require.Error(t, file.Close())
}

func TestFile_ErrorReadClosedFile(t *testing.T) {
	t.Parallel()

	// Create a file with some contents
	filePath := filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "no content")

	// Create a file File pointing to a non-existing file
	file, _ := linter.NewFile(filePath)

	// Close the File before reading
	require.NoError(t, file.Close())

	// Try to read the file
	_, err := file.Next()
	require.Error(t, err)
}
