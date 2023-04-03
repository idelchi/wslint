package linter_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/linter"
)

// Test the Formatter by
// 1. Creating a file with some contents
// 2. Creating a Formatter to interact with the file
// 3. Copying the file in reverse order to a replacement file
// 4. Replacing the original file with the replacement file
// 5. Checking that the contents are the same, as the reversed order.
func TestFormatter(t *testing.T) {
	t.Parallel()

	// Create a file with some contents
	filePath := filepath.Join(t.TempDir(), "test.txt")

	content := []string{
		"This is the original file.",
		"It has some dummy content,",
		"like this line.",
	}

	createFile(t, filePath, strings.Join(content, "\n"))

	// Create a file wrapper
	file, err := linter.NewFile(filePath)
	require.NoError(t, err)

	// Create a shadow
	require.NoError(t, err)

	formatter := linter.NewFormatter(file)

	contentReverse := []string{}

	// Read all lines into a slice
	content = iterateFile(t, file)

	// Copy in reverse order the content to the replacement file
	for i := len(content) - 1; i >= 0; i-- {
		contentReverse = append(contentReverse, content[i])
		require.NoError(t, formatter.Write(content[i]))
	}

	// Now replace the original file with the replacement
	require.NoError(t, formatter.Save())

	// Open the original file and check the lines
	require.NoError(t, file.Open())

	// Read all lines into a slice
	rows := iterateFile(t, file)

	// Check that the contents are the same, as the reversed order
	// Skip the last line as it is empty.
	require.Equal(t, contentReverse, rows[:len(rows)-1])
}

func TestFormatter_ErrorFolderNotExisting(t *testing.T) {
	t.Parallel()

	// Create a file wrapper, but pointing to a location that does not exist.
	// This should return an error from the CreateShadow function, which fails to create the 'shadow' file.
	filePath := filepath.Join(t.TempDir(), "no-existing-folder", "test.txt")
	file := &linter.File{Name: filePath}

	// Create a shadow
	_, err := linter.CreateShadow(file.Name)
	require.Error(t, err)
}
