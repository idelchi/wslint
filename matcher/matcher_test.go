// Tests for the matcher package.
package matcher_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/matcher"
)

// DummyLogger implements the Logger interface, but does nothing.
type DummyLogger struct{}

// Printf prints out to nothing.
func (DummyLogger) Printf(_ string, _ ...interface{}) {}

// Helper function to:
// - create a named file in a given folder
// - fill it with some content
// Returns the absolute path to the file.
func CreateTempFile(t *testing.T, dir, name, content string) (file string) {
	t.Helper()

	// Construct the path to the new file
	file = filepath.FromSlash(filepath.Join(dir, name))

	// Get the absolute path to the file.
	file, err := filepath.Abs(file)

	require.NoError(t, err)

	// Create the file.
	require.NoError(t, os.WriteFile(file, []byte(content), 0o600))

	return file
}

// TestGlobber_Match tests the Match function.
// All tests are run in a temporary folder, with a list of files created in it.
// It tests combinations of
// - text files
// - binary files
// - hidden files
// - multiple files.
func TestGlobber_Match(t *testing.T) {
	t.Parallel()

	logger := DummyLogger{}

	// Create a list of test cases.
	tcs := []struct {
		name     string   // Name of the test case (for logging)
		files    []string // List of files to create in the temporary folder
		content  []string // Content of each file
		expected []string // List of files that should be found by the globber
		hidden   bool     // Whether to allow the globber to match hidden files
		comment  string   // Comment in case of failure
	}{
		{
			comment:  "Expected to find the text file",
			name:     "Single text file",
			files:    []string{"test.txt"},
			content:  []string{"test"},
			expected: []string{"test.txt"},
			hidden:   false,
		},
		{
			comment:  "Expected to not find the single file, since it is empty and identified as binary",
			name:     "Single binary file",
			files:    []string{"test.txt"},
			content:  []string{""},
			expected: []string{},
		},
		{
			comment:  "Expected to find only a text file, since the others are either exe, binary or hidden",
			name:     "Multiple mixed files, text, exe, binary, hidden",
			files:    []string{"test.txt", "test.exe", "test", ".test.txt"},
			content:  []string{"test", "test", "", "test"},
			expected: []string{"test.txt"},
		},
		{
			comment:  "Expected to find the file, since hidden is set to true",
			name:     "Single hidden file",
			files:    []string{".test.txt"},
			content:  []string{"test"},
			expected: []string{".test.txt"},
			hidden:   true,
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a temporary folder
			// Each test needs its own folder, otherwise the globber will find files from other tests.
			dir := t.TempDir()

			// Append the dir name to the expected files.
			for i := range tc.expected {
				tc.expected[i] = filepath.ToSlash(filepath.Join(dir, tc.expected[i]))
			}

			// Create the files.
			for i := range tc.files {
				_ = CreateTempFile(t, dir, tc.files[i], tc.content[i])
			}

			// Create the globber.
			matcher := matcher.New(tc.hidden, []string{}, logger)

			// Match all files in the created directory.
			require.NoError(t, matcher.Match(dir+"/*"))

			// Get the list of files.
			files := matcher.ListFiles()

			require.Equal(t, len(tc.expected), len(files), "# files found failed: %s", tc.comment)
			require.Equal(t, tc.expected, files, "files found failed: %s", tc.comment)
		})
	}
}

// TestGlobber_Match_CornerCases tests the Match function with corner cases.
// It tests:
// - bad patterns
// - multiple globs
// The tests are collected in a single test, since they all operate on one folder & globber.
func TestGlobber_Match_CornerCases(t *testing.T) {
	t.Parallel()

	logger := DummyLogger{}

	// Create one file.
	dir := t.TempDir()
	file := CreateTempFile(t, dir, "test.txt", "test")

	// Create a list of test cases.
	tcs := []struct {
		name     string   // Name of the test case (for logging)
		globs    []string // List of globs to match
		expected int      // Number of files that should be found by the globber
		err      bool     // Whether an error is expected
		comment  string   // Comment in case of failure
	}{
		{
			comment: "Expected to get an error for a bad pattern",
			name:    "Bad pattern",
			globs:   []string{"["},
			err:     true,
		},
		{
			comment:  "Expected to find the file only once",
			name:     "Single file, multiple includes",
			globs:    []string{file, dir + "/*"},
			expected: 1,
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create the globber
			matcher := matcher.Globber{Logger: logger}

			fErrCheck := require.NoError
			if tc.err {
				fErrCheck = require.Error
			}

			for _, glob := range tc.globs {
				fErrCheck(t, matcher.Match(glob), "match failed: %s", tc.comment)
			}

			// Get the list of files.
			files := matcher.ListFiles()

			require.Len(t, files, tc.expected, "# files found failed: %s", tc.comment)
		})
	}
}
