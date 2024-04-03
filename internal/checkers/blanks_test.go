//go:build excluded

package checkers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/idelchi/wslint/internal/checkers"
)

// Test the Blanks struct.
// Test-sequence is:
// 1. Create a Blanks struct.
// 2. Call the Analyze method for each line.
// 3. Call the Finalize method.
// 4. Check the results.
func TestBlanks(t *testing.T) {
	t.Parallel()

	// Test cases for the Blanks struct.
	tcs := []struct {
		name    string   // Name of the test case (for logging)
		lines   []string // List of lines to check
		rows    []int    // List of rows that are blank (at the end)
		stop    int      // Stop row
		err     error    // Error that should be returned
		comment string   // Comment in case of failure
	}{
		{
			name: "no blank line",
			lines: []string{
				"This text sequence is a single line and ends with no blank line at the end.",
			},
			err:     checkers.ErrTooFewBlanks,
			comment: "Sequence with no blank line at the end.",
			rows:    []int{1},
		},
		{
			name: "no blank line, long text",
			lines: []string{
				"This text sequence ends with no blank line at the end.",
				"It has many rows, but no blank line at the end.",
				"It should return row 3, an error, and a stop indicated by 0.",
			},
			err:     checkers.ErrTooFewBlanks,
			comment: "Sequence with no blank line at the end",
			rows:    []int{3},
		},
		{
			name: "too many blank lines",
			lines: []string{
				"This text sequence ends with too many blank line at the end.",
				"",
				"",
			},
			err:     checkers.ErrTooManyBlanks,
			rows:    []int{2, 3},
			stop:    2,
			comment: "Sequence with several blank lines at the end.",
		},
		{
			name: "one blank line",
			lines: []string{
				"This text sequence ends with exactly one blank line at the end.",
				"",
			},
			rows:    []int{2},
			comment: "Sequence with exactly one blank line at the end.",
		},
		{
			name: "only one blank line",
			lines: []string{
				"",
			},
			rows:    []int{1},
			comment: "Sequence with exactly one blank line at the end.",
		},
		{
			name: "only blank lines",
			lines: []string{
				"",
				"",
			},
			err:     checkers.ErrTooManyBlanks,
			rows:    []int{1, 2},
			stop:    1,
			comment: "Sequence with several blank lines at the end.",
		},
		{
			name: "mixed too many blank lines",
			lines: []string{
				"This text sequence ends with too many blank line at the end.",
				"",
				"",
				"Even if there's text mass in between,",
				"",
				"",
				"it's still too many blank lines.",
				"",
				"",
			},
			err:     checkers.ErrTooManyBlanks,
			rows:    []int{8, 9},
			stop:    8,
			comment: "Sequence with several blank lines at the end.",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			blanks := checkers.Blanks{}

			for i, line := range tc.lines {
				blanks.Analyze(line, i+1)
			}

			blanks.Finalize()

			rows, err := blanks.Results()
			stop := blanks.Stop()

			assert.Equal(t, tc.rows, rows, "rows failed: %s", tc.comment)
			assert.Equal(t, tc.err, err, "errors failed: %s", tc.comment)
			assert.Equal(t, tc.stop, stop, "stops failed: %s", tc.comment)
		})
	}
}

// Test the Blanks Fix method.
// Test-sequence is:
// 1. Create a Blanks struct.
// 2. Call the Fix method for each line.
// 3. Check the results.
// Expected result is that the same line is returned, as the Fix method is a dummy.
func TestBlanks_Fix(t *testing.T) {
	t.Parallel()

	// Test cases for the Blanks Fix method.
	tcs := []struct {
		name    string // Name of the test case (for logging)
		line    string // Line to fix
		fixed   string // Fixed line
		comment string // Comment in case of failure
	}{
		{
			name:    "Blanks returns the same line",
			line:    "Blanks.Fix() method is a dummy.",
			fixed:   "Blanks.Fix() method is a dummy.",
			comment: "The exact same string has to be passed back.",
		},
		{
			name:    "Blanks always returns the same line",
			line:    "Blanks.Fix() method is a dummy, it returns the same string.",
			fixed:   "Blanks.Fix() method is a dummy, it returns the same string.",
			comment: "The exact same string has to be passed back.",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			blanks := checkers.Blanks{}

			assert.Equal(t, tc.fixed, blanks.Fix(tc.line), "fix failed: %s", tc.comment)
		})
	}
}
