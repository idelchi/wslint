package linter_test

import (
	"os"
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

func TestFile_Rename(t *testing.T) {
	t.Parallel()

	// Create a file with no content
	filePath := filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "")
	file, err := linter.NewFile(filePath)

	// No error should occur when opening
	require.NoError(t, err)

	require.True(t, file.Exists())

	// Create a new file path
	filePathNew := filepath.Join(t.TempDir(), "renamed.txt")

	// Rename the file
	require.NoError(t, file.Rename(filePathNew))

	// Check that the old file doesn't exist using os.Stat
	_, err = os.Stat(filePath)

	require.ErrorIs(t, err, os.ErrNotExist)

	// Check that the file exists
	require.True(t, file.Exists())

	// Check that the file has been renamed
	require.Equal(t, filePathNew, file.Name)

	// Save it
	require.NoError(t, file.Save())
}

func TestFile_CloseClosedFile(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "")

	file, err := linter.NewFile(filePath)

	require.NoError(t, err)

	require.NoError(t, file.Close())
	require.NoError(t, file.Close())
}

// Test the File by
// 1. Creating an empty folder
// 2. Creating a File to interact with a non-existing file
// 3. Checking that an error occurs when opening the file.
func TestFile_ErrorFileNotExisting(t *testing.T) {
	t.Parallel()

	// Non-existing file
	filePath := filepath.Join(t.TempDir(), "test.txt")
	_, err := linter.NewFile(filePath)

	// Require an error when opening a non-existing file
	require.Error(t, err)
}

func TestFile_ErrorOperateUnopenedFile(t *testing.T) {
	t.Parallel()

	// Non-existing file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Create a File pointing to a non-existing file
	file, err := linter.NewFile(filePath)

	require.Error(t, err)

	// Close
	require.Error(t, file.Close())
	// Rename
	require.Error(t, file.Rename(filepath.Join(tmpDir, "test-renamed.txt")))
	// Replace
	require.Error(t, file.ReplaceWith(&linter.File{}))
}

func TestFile_ErrorReadClosedFile(t *testing.T) {
	t.Parallel()

	// Non-existing file
	filePath := filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "")

	// Create a File pointing to a non-existing file
	file, err := linter.NewFile(filePath)

	require.NoError(t, err)
	require.NoError(t, file.Close())

	// Read
	_, err = file.Next()
	require.Error(t, err)
}
