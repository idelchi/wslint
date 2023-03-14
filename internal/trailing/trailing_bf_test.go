package trailing_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/trailing"
)

// Benchmarks were used to compare the performance of different implementations of the Has function.
func benchmarkHas(b *testing.B, s string) {
	b.Helper()

	for n := 0; n < b.N; n++ {
		trailing.Has(s)
	}
}

// Benchmarks were used to compare the performance of different implementations of the Trim function.
func benchmarkTrim(b *testing.B, s string) {
	b.Helper()

	for n := 0; n < b.N; n++ {
		trailing.Trim(s)
	}
}

func BenchmarkTrailing(b *testing.B) {
	tcs := []struct {
		name string // Name of the test case (for logging)
		line string // Line to check
	}{
		{
			name: "empty",
		},
		{
			name: "long",
			line: "xxxxxxxxxxxxxxxxx",
		},
		{
			name: "whitespace",
			line: "xxxxxxxxxxxxxxxxx  ",
		},
		{
			name: "whitespaces",
			line: "xxxxxxxxxxxxxxxxx                                                                          ",
		},
	}

	for _, tc := range tcs {
		name := "(Has): " + tc.name
		b.Run(name, func(b *testing.B) {
			benchmarkHas(b, tc.line)
		})
	}

	for _, tc := range tcs {
		name := "(Trim): " + tc.name
		b.Run(name, func(b *testing.B) {
			benchmarkTrim(b, tc.line)
		})
	}
}

// FuzzHas is not the best example of a fuzzing test, but it's here for learning purposes.
// It brought some insight into how the unicode.IsSpace function works in combination with the strings.TrimRightFunc
// function.
// TODO(Idelchi): Follow the example [here]: https://go.dev/doc/tutorial/fuzz
// It explains bytes, rune, non-utf8 etc.
func FuzzHas(f *testing.F) {
	f.Add("A line with no whitespace.")
	f.Add("A line with whitespace.  ")
	f.Add("A line with whitespace tabs.\t")

	f.Fuzz(func(t *testing.T, line string) {
		if trailing.Has(line) {
			require.NotEqual(t, line, trailing.Trim(line))
		} else {
			require.Equal(t, line, trailing.Trim(line))
		}
	})
}
