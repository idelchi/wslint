//go:build excluded

package linter_test

import (
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
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			file := CreateTempFile(t, tc.content...)

			lintFile := linter.New(file)

			reader, writer := NewReaderWriter(t, file)

			require.NoError(t, lintFile.Lint(reader))

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
			require.NoError(t, lintFile.Fix(writer))
			lintFile.Summary()

			// Lint it again
			lintFile = linter.New(file)
			require.NoError(t, lintFile.Lint(writer))

			for _, c := range lintFile.Checkers {
				_, err := c.Results()
				require.NoError(t, err)
			}
			lintFile.Summary()
		})
	}
}

// NewReaderFormatter creates a new reader and formatter for the given file.
func NewReaderWriter(t *testing.T, file string) (*linter.Reader, *linter.Writer) {
	t.Helper()

	reader, err := linter.NewReader(file)
	require.NoError(t, err)

	writer, err := linter.NewWriter(reader)
	require.NoError(t, err)

	return reader, writer
}

func TestLinter_InsertChecker(t *testing.T) {
	t.Parallel()

	lintFile := linter.Linter{}
	lintFile.InsertChecker(&checkers.Blanks{})
	require.Len(t, lintFile.Checkers, 1)
}

func TestLinter_ErrorNoCheckersConfigured(t *testing.T) {
	t.Parallel()

	// Create the file
	lintFile := linter.Linter{Name: CreateTempFile(t, "This file ends with no whitespace.")}

	reader, writer := NewReaderWriter(t, lintFile.Name)

	require.Error(t, lintFile.Lint(reader))
	require.Error(t, lintFile.Fix(writer))
}

func TestLinter_ErrorSummary(t *testing.T) {
	t.Parallel()

	lintFile := linter.Linter{}
	lintFile.Error = assert.AnError
	require.False(t, lintFile.Summary())
}
