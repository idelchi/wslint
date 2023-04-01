package linter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/linter"
)

func TestFileReplacer(t *testing.T) {
	t.Parallel()

	// Create a file with some contents
	filePath := filepath.Join(t.TempDir(), "test.txt")

	content := []string{
		"This is the original file.",
		"It has some dummy content,",
		"like this line.",
	}

	createFile(t, filePath, strings.Join(content, "\n"))

	// Create a file manager
	manager := linter.Manager{Name: filePath}
	require.NoError(t, manager.Open())

	// Create a replacer
	replacer := linter.Replacer{Original: manager}
	require.NoError(t, replacer.Setup())

	contentReverse := []string{}

	// Copy in reverse order the content to the replacement file
	for i := len(content) - 1; i >= 0; i-- {
		contentReverse = append(contentReverse, content[i])
		require.NoError(t, replacer.Write(content[i]))
	}

	// Now replace the original file with the replacement
	require.NoError(t, replacer.Replace())

	// Open the original file and check the lines
	require.NoError(t, manager.Open())

	// Read all lines into a slice
	rows := iterateManager(t, manager)

	// Check that the contents are the same, as the reversed order
	// Skip the last line as it is empty.
	require.Equal(t, contentReverse, rows[:len(rows)-1])
}

// TODO(Idelchi): Table driven tests.
func TestFileReplacer_Error(t *testing.T) {
	t.Parallel()

	// 1. Try to setup in a folder that does not exist.
	filePath := filepath.Join(t.TempDir(), "no-existing-folder", "test.txt")
	replacer := linter.Replacer{Original: linter.Manager{Name: filePath}}
	require.Error(t, replacer.Setup())

	// 2. Fail to Copy a line
	// Create a file with some contents
	filePath = filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "content")
	replacer = linter.Replacer{Original: linter.Manager{Name: filePath}}
	require.NoError(t, replacer.Setup())
	// Close both files
	require.NoError(t, replacer.Original.Open())
	require.NoError(t, replacer.Close())
	require.Error(t, replacer.Write("This is a line which should fail"))

	// Start over
	require.NoError(t, replacer.Setup())
	// Delete the temporary file
	require.NoError(t, os.Remove(replacer.Replacement.Name))
	// Try to replace the original file
	require.Error(t, replacer.Replace())

	// Start over
	replacer = linter.Replacer{Original: linter.Manager{Name: filePath}}
	require.NoError(t, replacer.Setup())
	// Never open the original file
	// Try to replace the original file
	require.Error(t, replacer.Replace())

	// Start over
	replacer = linter.Replacer{Original: linter.Manager{Name: filePath}}
	require.NoError(t, replacer.Setup())
	// Never open the original file
	// Try to close the original file
	var errPtr error

	require.Error(t, replacer.Close(&errPtr))
	require.Error(t, errPtr)
	require.Error(t, replacer.Close(&errPtr))
	require.Error(t, errPtr)

	// Start over
	replacer = linter.Replacer{Original: linter.Manager{Name: filePath}}
	require.NoError(t, replacer.Original.Open())
	// Never setup the replacer
	// Try to close the original file
	require.Error(t, replacer.Close())

	// Start over
	replacer = linter.Replacer{Original: linter.Manager{Name: filePath}}
	// Never open the original file
	// Never setup the replacer
	// Try to close the files
	require.Error(t, replacer.Close())
}