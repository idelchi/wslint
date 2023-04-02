package linter_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/checkers"
	"github.com/idelchi/wslint/internal/linter"
)

func TestLinter(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string   // Name of the test case (for logging)
		content []string // Content of file to generate
		errs    []error  // Errors that should be returned
	}{
		{
			name: "OK",
			content: []string{
				"This file ends with no whitespace and a blank line at the end.",
				"",
			},
		},
		{
			name: "Missing blank line",
			content: []string{
				"This file ends with no whitespace but misses a blank line at the end.",
			},
			errs: []error{nil, checkers.ErrTooFewBlanks},
		},
		{
			name: "Many blank lines",
			content: []string{
				"This file ends with no whitespace but too many blank lines at the end.",
				"",
				"",
			},
			errs: []error{nil, checkers.ErrTooManyBlanks},
		},
		{
			name: "All blank lines",
			content: []string{
				"",
				"",
			},
			errs: []error{nil, checkers.ErrTooManyBlanks},
		},
		{
			name: "Trailing whitespace",
			content: []string{
				"This file ends with trailing whitespace but a blank line at the end. ",
				"",
			},
			errs: []error{checkers.ErrHasTrailing, nil},
		},
		{
			name: "Whitespace and blanks issue",
			content: []string{
				"This file ends with whitespace and has too many blank lines at the end. ",
				"",
				"",
			},
			errs: []error{checkers.ErrHasTrailing, checkers.ErrTooManyBlanks},
		},
		{
			name: "Mixed whitespace",
			content: []string{
				"This file ends with mixed whitespace but a blank line at the end. \t",
				"",
			},
			errs: []error{checkers.ErrHasTrailing, nil},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			file := filepath.Join(t.TempDir(), "test.txt")

			createFile(t, file, strings.Join(tc.content, "\n"))

			lintFile, err := linter.NewLinter(file)

			require.NoError(t, err)

			require.NoError(t, lintFile.Lint())

			for i, c := range lintFile.Checkers {
				_, err := c.Results()
				if tc.errs != nil {
					require.Equal(t, tc.errs[i], err)
				} else {
					require.NoError(t, err)
				}
			}
			lintFile.Summary()

			// Fix the file
			require.NoError(t, lintFile.Fix())
			lintFile.Summary()

			// Lint it again
			lintFile, err = linter.NewLinter(file)
			require.NoError(t, err)
			require.NoError(t, lintFile.Lint())

			for _, c := range lintFile.Checkers {
				_, err := c.Results()
				require.NoError(t, err)
			}
			lintFile.Summary()
		})
	}
}

func TestLinter_InsertChecker(t *testing.T) {
	t.Parallel()

	lintFile := linter.Linter{}
	lintFile.InsertChecker(&checkers.Blanks{})
	require.Len(t, lintFile.Checkers, 1)
}

func TestLinter_ErrorFileNotExisting(t *testing.T) {
	t.Parallel()

	lintFile, err := linter.NewLinter(filepath.Join(t.TempDir(), "test.txt"))
	require.Error(t, err)
	require.Error(t, lintFile.Fix())
	lintFile.Summary()
}

func TestLinter_ErrorNoCheckersConfigured(t *testing.T) {
	t.Parallel()

	// Create the file
	filePath := filepath.Join(t.TempDir(), "test.txt")
	createFile(t, filePath, "This file ends with no whitespace.")
	lintFile := linter.Linter{Name: filePath}
	require.Error(t, lintFile.Fix())
}

func TestLinter_ErrorSummary(t *testing.T) {
	t.Parallel()

	lintFile := linter.Linter{}
	lintFile.Error = assert.AnError
	require.False(t, lintFile.Summary())
}
