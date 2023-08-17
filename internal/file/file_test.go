//go:build excluded

package linter_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/linter"
)

// Test the Reader by
// 1. Creating a file with some contents
// 2. Creating a Reader to interact with the file
// 3. Reading all lines
// 4. Checking that the contents in steps (1) and (3) are the same.
func TestReader(t *testing.T) {
	t.Parallel()

	content := []string{
		"Line 1",
		"Line 2",
		"Line 3",
	}

	// Create a file with some contents and hand it to the Reader
	file, err := linter.NewReader(CreateTempFile(t, content...))

	// No error should occur when opening
	require.NoError(t, err)

	// Read all lines into a slice
	rows := ReadAll(t, file)

	// Check that the contents are the same
	require.Equal(t, content, rows)
}

// Test the Write method of Reader by
// 1. Creating a file with no content
// 2. Creating a Reader to interact with the file
// 3. Writing some lines to the file
// 4. Reading all lines
// 5. Checking that the contents in steps (3) and (4) are the same.
func TestReader_Write(t *testing.T) {
	t.Parallel()

	// Create a file and hand it to the Reader
	file, err := linter.NewReader(CreateTempFile(t))

	// No error should occur when opening
	require.NoError(t, err)

	content := []string{
		"Line 1",
		"Line 2",
		"Line 3",
	}

	// Write some lines to the file
	// No error should occur
	require.NoError(t, file.Write(content...))

	// Read all lines into a slice
	rows := ReadAll(t, file)

	// Check that the contents are the same
	// The last line is empty because the file ends with a newline
	require.Equal(t, content, rows[:len(rows)-1])
}

// Test the Rename method of Reader by
// 1. Creating a file
// 2. Creating a Reader to interact with the file
// 3. Renaming the file
// 4. Checking that the old file doesn't exist
// 5. Checking that the renamed file exists
// 6. Saving the file.
func TestReader_Rename(t *testing.T) {
	t.Parallel()

	// Create a file and hand it to the Reader
	filePath := CreateTempFile(t)
	file, err := linter.NewReader(filePath)

	// No error should occur when opening
	require.NoError(t, err)

	// Create a new file path and rename the file
	require.NoError(t, file.Rename(TempFileName(t)))

	// Check that the old file doesn't exist
	require.NoFileExists(t, filePath)

	// Check that the renamed file exists
	require.FileExists(t, file.Name)

	// Save it
	require.NoError(t, file.Save())
}

// Test the Replace method of Reader by
// 1. Creating a file
// 2. Creating a Reader to interact with the file
// 3. Renaming the file
// 4. Checking that the old file doesn't exist
// 5. Checking that the renamed file exists
// 6. Saving the file.
func TestReader_Replace(t *testing.T) {
	t.Parallel()

	// Create a file and hand it to the Reader
	filePath1 := CreateTempFile(t)
	file1, err := linter.NewReader(filePath1)
	require.NoError(t, err)

	// Create another file and hand it to the Reader
	content := []string{
		"Replacement test",
	}
	filePath2 := CreateTempFile(t, content...)
	file2, err := linter.NewReader(filePath2)
	require.NoError(t, err)

	// Replace the first file with the second
	require.NoError(t, file1.ReplaceWith(file2))

	// Check that the replacement file doesn't exist
	require.NoFileExists(t, filePath2)

	// Reload the file
	require.NoError(t, file1.Load(filePath1))
	rows := ReadAll(t, file1)

	// Check that the contents are the same
	require.Equal(t, content, rows)
}

// Test the Close method of Reader by
// 1. Creating a file
// 2. Creating a Reader to interact with the file
// 3. Closing the file multiple times, with the expectation that no error occurs.
func TestReader_Close(t *testing.T) {
	t.Parallel()

	// Create a file with some contents and hand it to the Reader
	file, err := linter.NewReader(CreateTempFile(t))

	t.Cleanup(func() {
		// Close the file
		require.NoError(t, file.Close())
	})

	// No error should occur when opening
	require.NoError(t, err)

	// Close the file n times
	for i := 0; i < 10; i++ {
		require.NoError(t, file.Close())
	}
}

// Test the error case of interacting with a non-existing file
// 1. Creating an empty folder
// 2. Creating a Reader to interact with a non-existing file
// 3. Checking that an error occurs.
func TestReader_ErrorFileNotExisting(t *testing.T) {
	t.Parallel()

	// Create a path to a non-existing file and hand it to the Reader
	filePath := TempFileName(t)
	_, err := linter.NewReader(filePath)

	// Require an error when opening a non-existing file
	require.Error(t, err)
}

func TestReader_ErrorOperateUnopenedFile(t *testing.T) {
	t.Parallel()

	// Non-existing file
	filePath := TempFileName(t)

	// Create a File pointing to a non-existing file
	file, err := linter.NewReader(filePath)

	require.Error(t, err)

	// Reset
	require.Error(t, file.Reset())
	// Write
	require.Error(t, file.Write("test"))
	// Close
	require.Error(t, file.Close())
	// Rename
	require.Error(t, file.Rename(filepath.Join(t.TempDir(), "test-renamed.txt")))
	// Replace
	require.Error(t, file.ReplaceWith(&linter.Reader{}))
}

func TestReader_ErrorReadClosedFile(t *testing.T) {
	t.Parallel()

	file, err := linter.NewReader(CreateTempFile(t))

	require.NoError(t, err)
	require.NoError(t, file.Close())

	// Read
	_, err = file.Next()
	require.Error(t, err)
}

func TestReader_ErrorRename(t *testing.T) {
	t.Parallel()

	file, err := linter.NewReader(CreateTempFile(t))

	require.NoError(t, err)

	// Create a temp path to a non-existing folder
	tempPath := filepath.Join(t.TempDir(), "does-not-exist", "test.txt")

	require.Error(t, file.Rename(tempPath))
}
