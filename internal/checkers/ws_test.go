//go:build excluded

package checkers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/checkers"
)

// Test the Whitespace struct.
// Test-sequence is:
// 1. Create a Whitespace struct.
// 2. Call the Analyze method for each line.
// 3. Call the Finalize method.
// 4. Check the results.
func TestWhiteSpace(t *testing.T) {
	t.Parallel()

	// Test cases for the Whitespace struct.
	tcs := []struct {
		name    string   // Name of the test case (for logging)
		lines   []string // List of lines to check
		rows    []int    // List of rows that should have trailing whitespace
		stop    int      // Stop row
		err     error    // Error that should be returned
		comment string   // Comment in case of failure
	}{
		{
			name: "no trailing whitespace",
			lines: []string{
				"This line has no trailing whitespace.",
				"And neither does this.",
			},
			comment: "Sequence with no trailing whitespace.",
		},
		{
			name: "Two rows with trailing whitespace",
			lines: []string{
				"This line has no trailing whitespace.",
				"But this one does this. ",
				"And so does this one. \t",
				"Here too. \t    ",
			},
			rows:    []int{2, 3, 4},
			err:     checkers.ErrHasTrailing,
			comment: "Sequence with trailing whitespace.",
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			linter := checkers.Whitespace{}

			for row, line := range tc.lines {
				linter.Analyze(line, row+1)
			}

			linter.Finalize()

			rows, err := linter.Results()
			stop := linter.Stop()

			require.Equal(t, tc.rows, rows, "rows failed: %s", tc.comment)
			require.Equal(t, tc.err, err, "err failed: %s", tc.comment)
			require.Equal(t, tc.stop, stop, "stop failed: %s", tc.comment)
		})
	}
}

func TestWhiteSpace_Fix(t *testing.T) {
	t.Parallel()

	// Test cases for the whitespace.Fix() method.
	tcs := []struct {
		name    string // Name of the test case (for logging)
		line    string // Line to check
		fixed   string // Line after whitespace is removed
		comment string // Comment in case of failure
	}{
		{
			name:    "no trailing whitespace",
			line:    "This line has no trailing whitespace.",
			fixed:   "This line has no trailing whitespace.",
			comment: "Sequence with no trailing whitespace.",
		},
		{
			name:    "Trailing space",
			line:    "This line has trailing space.   ",
			fixed:   "This line has trailing space.",
			comment: "Sequence with trailing spaces.",
		},
		{
			name:    "Trailing tabs",
			line:    "This line has trailing tabs.\t\t",
			fixed:   "This line has trailing tabs.",
			comment: "Sequence with trailing tabs.",
		},
		{
			name:    "Trailing mix",
			line:    "This line has mixed trailing tabs and spaces.\t  \t  ",
			fixed:   "This line has mixed trailing tabs and spaces.",
			comment: "Sequence with trailing tabs and spaces.",
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			linter := checkers.Whitespace{}

			require.Equal(t, tc.fixed, linter.Fix(tc.line), "fix failed: %s", tc.comment)
		})
	}
}
