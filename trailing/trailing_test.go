// Tests "Has" and "Trim" functions with a sequence of test cases, mixing empty lines, trailing spaces and trailing
// tabs.

package trailing_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/trailing"
)

func TestTrailing(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string // Name of the test case (for logging)
		line    string // Line to check
		has     bool   // Whether the line has trailing whitespace
		trimmed string // The line with trailing whitespace (manually) removed
	}{
		{
			name: "empty line",
		},
		{
			name: "trailing line",
			line: " ",
			has:  true,
		},
		{
			name:    "trailing spaces",
			line:    "test ",
			has:     true,
			trimmed: "test",
		},
		{
			name:    "trailing tabs",
			line:    "test\t",
			has:     true,
			trimmed: "test",
		},
		{
			name:    "trailing spaces and tabs",
			line:    "test \t",
			has:     true,
			trimmed: "test",
		},
		{
			name:    "trailing tabs and space",
			line:    "test\t ",
			has:     true,
			trimmed: "test",
		},
		{
			name:    "trailing spaces and tabs with text",
			line:    "test \ttext",
			trimmed: "test \ttext",
		},
		{
			name:    "trailing spaces and tabs with text and spaces",
			line:    "test \ttext  ",
			has:     true,
			trimmed: "test \ttext",
		},
		{
			name:    "trailing spaces and tabs with text and tabs",
			line:    "test \ttext\t",
			has:     true,
			trimmed: "test \ttext",
		},
		{
			name:    "trailing spaces and tabs with text and spaces and tabs",
			line:    "test \ttext  \t\t",
			has:     true,
			trimmed: "test \ttext",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.has, trailing.Has(tc.line), "Has() failed: %q", tc.line)
			require.Equal(t, tc.trimmed, trailing.Trim(tc.line), "Trim() failed: %q", tc.line)
		})
	}
}
