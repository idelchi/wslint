package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContains tests the contains function.
func TestContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		files    []string
		expected bool
	}{
		{
			name:     "empty list",
			file:     "foo",
			files:    []string{},
			expected: false,
		},
		{
			name:     "not present",
			file:     "foo",
			files:    []string{"bar", "baz"},
			expected: false,
		},
		{
			name:     "present",
			file:     "foo",
			files:    []string{"foo", "bar", "baz"},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, contains(test.file, test.files))
		})
	}
}

// TestIsExplicitlyIncluded tests the isExplicitlyIncluded function.
func TestIsExplicitlyIncluded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		expected bool
	}{
		{
			name:     "empty",
			file:     "",
			expected: true,
		},
		{
			name:     "no globs",
			file:     "foo",
			expected: true,
		},
		{
			name:     "glob only in filename",
			file:     "foo/*.txt",
			expected: true,
		},
		{
			name:     "globs in path",
			file:     "foo/*/bar",
			expected: false,
		},
		{
			name:     "globs in name",
			file:     "foo*",
			expected: false,
		},
		{
			name:     "globs in extension",
			file:     "foo.*",
			expected: false,
		},
		{
			name:     "globs in name and extension",
			file:     "foo*.*",
			expected: false,
		},
		{
			name:     "globs in path and name",
			file:     "foo/*/bar*",
			expected: false,
		},
		{
			name:     "globs in path and extension",
			file:     "foo/*/bar.*",
			expected: false,
		},
		{
			name:     "globs in path and name and extension",
			file:     "foo/*/bar*.*",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, IsExplicitlyIncluded(test.file))
		})
	}
}
