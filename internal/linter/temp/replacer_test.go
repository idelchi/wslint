package linter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/linter"
)

// Test the Replacer by
// 1. Creating a file with some contents
// 2. Creating a Reader to interact with the file
// 3. Creating a Replacer to interact with the file
// 4. Copying the file in reverse order to a replacement file
// 5. Replacing the original file with the replacement file
// 6. Checking that the contents are the same, as the reversed order.
func TestReplacer(t *testing.T) {
	t.Parallel()

	// Create a file with some contents
	filePath := filepath.Join(t.TempDir(), "test.txt")

	content := []string{
		"This is the original file.",
		"It has some dummy content,",
		"like this line.",
	}

	createFile(t, filePath, strings.Join(content, "\n"))

	// Create a reader
	reader, err := linter.NewReader(filePath)
	require.NoError(t, err)

	// Create a replacer
	replacer, err := linter.NewReplacer(reader)
	require.NoError(t, err)

	contentReverse := []string{}

	// Copy in reverse order the content to the replacement file
	for i := len(content) - 1; i >= 0; i-- {
		contentReverse = append(contentReverse, content[i])
		require.NoError(t, replacer.Write(content[i]))
	}

	// Now replace the original file with the replacement
	require.NoError(t, replacer.Replace())

	// Open the original file and check the lines
	require.NoError(t, reader.Open())

	// Read all lines into a slice
	rows := iterateReader(t, reader)

	// Check that the contents are the same, as the reversed order
	// Skip the last line as it is empty.
	require.Equal(t, contentReverse, rows[:len(rows)-1])
}

func TestReplacer_ErrorFolderNotExisting(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "no-existing-folder", "test.txt")
	_, err := linter.NewReplacer(linter.Reader{Name: filePath})
	require.Error(t, err)
}

func TestReplacer_ErrorCopyToClosedFile(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "test.txt")
	reader := linter.Reader{Name: filePath}
	replacer, err := linter.NewReplacer(reader)

	require.NoError(t, err)
	require.NoError(t, replacer.Replacement.Close())

	require.Error(t, replacer.Write("This is a line which should fail"))
}

func TestReplacer_ErrorReplaceWhenDeleted(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "")
	reader, _ := linter.NewReader(filePath)
	replacer, err := linter.NewReplacer(reader)

	require.NoError(t, err)

	require.NoError(t, os.Remove(replacer.Replacement.Name))

	require.Error(t, replacer.Replace())
}

func TestReplacer_ErrorReplaceUnOpened(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "test.txt")
	reader := linter.Reader{Name: filePath}
	replacer, err := linter.NewReplacer(reader)

	require.NoError(t, err)

	require.Error(t, replacer.Replace())
}

func TestReplacer_ErrorCloseUnOpened(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "test.txt")
	reader := linter.Reader{Name: filePath}
	replacer, err := linter.NewReplacer(reader)

	require.NoError(t, err)

	var errPtr error

	require.Error(t, replacer.Close(&errPtr))
	require.Error(t, errPtr)
	require.Error(t, replacer.Close(&errPtr))
	require.Error(t, errPtr)
}

func TestReplacer_ErrorCloseNotSetup(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "test.txt")

	createFile(t, filePath, "")

	reader, err := linter.NewReader(filePath)
	require.NoError(t, err)

	replacer := linter.Replacer{Original: reader}
	require.Error(t, replacer.Close())
}

func TestReplacer_ErrorCloseUnopenedAndNotSetup(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(t.TempDir(), "test.txt")

	reader := linter.Reader{Name: filePath}

	replacer := linter.Replacer{Original: reader}
	require.Error(t, replacer.Close())
}
