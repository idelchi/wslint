package linter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/idelchi/wslint/internal/checkers"
	"github.com/idelchi/wslint/internal/linter"

	"bou.ke/monkey"
)

// Create a file in a temporary folder, fill it with some content, and close it.
func createFile(t *testing.T, file, content string) {
	t.Helper()

	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

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

			lintFile := linter.New(file)

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
			lintFile = linter.New(file)
			require.NoError(t, lintFile.Lint())

			for _, c := range lintFile.Checkers {
				_, err := c.Results()
				require.NoError(t, err)
			}
			lintFile.Summary()
		})
	}
}

func TestLinter_Error(t *testing.T) {
	t.Parallel()

	// File doesn't exist
	file := filepath.Join(t.TempDir(), "test.txt")
	lintFile := linter.Linter{Name: file}
	lintFile.InsertChecker(&checkers.Blanks{})
	require.Error(t, lintFile.Lint())
	lintFile.Summary()
	require.Error(t, lintFile.Fix())
	lintFile.Summary()

	// Create the file
	createFile(t, file, "This file ends with no whitespace. It however has some illegal characters.")
	// Remove the checkers
	lintFile.Checkers = nil
	require.Error(t, lintFile.Lint())
	lintFile.Summary()
	require.Error(t, lintFile.Fix())
	lintFile.Summary()

	lintFile.Error = assert.AnError
	require.False(t, lintFile.Summary())
}

// TODO(Idelchi): Replace with a proper test case, not monkey-patching. It's an indicator that the code needs to be
// refactored.
func TestLinter_Error_MonkeyPatch(t *testing.T) { //nolint:paralleltest // Monkey patching is not thread-safe
	// Make life easier - assume there's an environment variable called 'RUN_MONKEY_PATCH_TESTS'
	if ok, err := strconv.ParseBool(os.Getenv("RUN_MONKEY_PATCH_TESTS")); !ok || err != nil {
		t.Skip("Skipping test because it requires monkey patching (inlining disabled)")
	}

	file := filepath.Join(t.TempDir(), "test.txt")
	lintFile := &linter.Linter{Name: file}
	lintFile.InsertChecker(&checkers.Blanks{})
	// Create the file
	createFile(t, file, "This file ends with whitespace.  ")

	// Monkey Patching!
	defer monkey.UnpatchAll()

	var manager *linter.Manager

	// Next in Lint
	monkey.PatchInstanceMethod(reflect.TypeOf(manager), "Next", func(*linter.Manager) (string, error) {
		return "", fmt.Errorf("failed to read line: %w", assert.AnError)
	})
	require.Error(t, lintFile.Lint())
	monkey.UnpatchAll()

	// Open in Fix
	monkey.PatchInstanceMethod(reflect.TypeOf(manager), "Open", func(*linter.Manager) error {
		return fmt.Errorf("failed to open file: %w", assert.AnError)
	})

	monkey.PatchInstanceMethod(reflect.TypeOf(lintFile), "Lint", func(*linter.Linter) error {
		return nil
	})

	require.Error(t, lintFile.Fix())

	monkey.UnpatchInstanceMethod(reflect.TypeOf(manager), "Open")

	// Setup in Fix
	var replacer *linter.Replacer

	monkey.PatchInstanceMethod(reflect.TypeOf(replacer), "Setup", func(*linter.Replacer) error {
		return fmt.Errorf("failed to setup file: %w", assert.AnError)
	})

	require.Error(t, lintFile.Fix())

	monkey.UnpatchInstanceMethod(reflect.TypeOf(replacer), "Setup")

	// Next in Fix
	monkey.PatchInstanceMethod(reflect.TypeOf(manager), "Next", func(*linter.Manager) (string, error) {
		return "", fmt.Errorf("failed to read line: %w", assert.AnError)
	})

	require.Error(t, lintFile.Fix())

	monkey.UnpatchInstanceMethod(reflect.TypeOf(manager), "Next")

	// Write in Fix
	monkey.PatchInstanceMethod(reflect.TypeOf(replacer), "Write", func(*linter.Replacer, string) error {
		return fmt.Errorf("failed to write file: %w", assert.AnError)
	})

	require.Error(t, lintFile.Fix())

	monkey.UnpatchInstanceMethod(reflect.TypeOf(replacer), "Write")

	// Replace in Fix
	monkey.PatchInstanceMethod(reflect.TypeOf(replacer), "Replace", func(*linter.Replacer) error {
		return fmt.Errorf("failed to replace file: %w", assert.AnError)
	})

	require.Error(t, lintFile.Fix())

	monkey.UnpatchInstanceMethod(reflect.TypeOf(replacer), "Replace")
}
