package stuttering_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/pkg/stuttering"
)

func TestStuttering(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string   // Name of the test case (for logging)
		line     string   // Line to check
		has      bool     // Whether the line has stuttering words
		stutters []string // The stuttering words identified
		trimmed  string   // The line with stuttering words (first occurrence) removed
	}{
		{
			name: "no words",
			line: "",
		},
		{
			name:    "single word",
			line:    "hello",
			trimmed: "hello",
		},
		{
			name:     "stuttering pair",
			line:     "hello hello", //nolint:dupword // This is a stuttering pair for testing purposes.
			has:      true,
			stutters: []string{"(hello hello)"},
			trimmed:  "hello",
		},
		{
			name:     "stuttering pair with punctuation",
			line:     "hello hello!",
			has:      true,
			stutters: []string{"(hello hello!)"},
			trimmed:  "hello!",
		},
		{
			name:     "multiple stuttering pairs",
			line:     "hey hey! hello hello! hi hi!",
			has:      true,
			stutters: []string{"(hey hey!)", "(hello hello!)", "(hi hi!)"},
			trimmed:  "hey! hello! hi!",
		},
		{
			name:    "non-stuttering pairs",
			line:    "hey hello hi",
			trimmed: "hey hello hi",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.has, stuttering.Has(tc.line), "Has() failed: %q", tc.line)
			require.ElementsMatch(t, tc.stutters, stuttering.Find(tc.line), "Find() failed: %q", tc.line)
			require.Equal(t, tc.trimmed, stuttering.Trim(tc.line), "Trim() failed: %q", tc.line)
		})
	}
}
