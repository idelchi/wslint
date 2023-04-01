// Package matcher provides a utility for matching files based on glob patterns, with support for
// excluding directories, hidden files and folders, and binary files.
package matcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/godoc/util"
	"golang.org/x/tools/godoc/vfs"

	"github.com/bmatcuk/doublestar/v4"
)

// Logger is an interface for logging formatted messages.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Globber is a file matcher that compiles a list of files matching a given pattern, while
// excluding files based on provided exclude patterns and options.
type Globber struct {
	// Exclude is a list of patterns that are used to exclude files.
	Exclude []string
	// Logger is a logger for debug messages (mainly).
	Logger Logger
	// files is the list of files that are added to the matcher, after matching and applying the options.
	files []string
}

// ListFiles lists all files found by the Globber.
func (m *Globber) ListFiles() []string {
	if m.files == nil {
		return []string{}
	}

	return m.files
}

// New creates a Globber with default settings, including:
//   - Excluding the executable itself
//   - Excluding all kinds of executables
//   - Excluding specific known directories
//   - Excluding hidden folders & files if the 'hidden' parameter is set to false
//
// The 'exclude' parameter allows specifying additional exclude patterns, and the 'logger' parameter
// sets the logger for debug messages.
func New(hidden bool, exclude []string, logger Logger) Globber {
	matcher := Globber{
		Exclude: exclude,
		Logger:  logger,
	}

	// Get the name of the executable itself
	if exe, err := os.Executable(); err == nil {
		// Exclude the executable itself
		matcher.Exclude = append(matcher.Exclude, exe)
	}

	matcher.Exclude = append(matcher.Exclude,
		// Exclude all kinds of executables
		"**/*.exe",

		// Exclude some known directories
		"**/.git/**",
		"**/.vscode-server/**",
		"**/node_modules/**",
		"**/vendor/**",
		"**/.task/**",
		"**/.cache/**",
	)

	if !hidden {
		// Exclude hidden folders & files if hidden is false
		matcher.Exclude = append(matcher.Exclude, "**/.*", "**/.*/**/*")
	}

	return matcher
}

// isBinary returns true if the given file is detected as a binary file, false otherwise.
func (m *Globber) isBinary(file string) bool {
	fs := vfs.OS(filepath.Dir(file))

	return !util.IsTextFile(fs, filepath.Base(file))
}

// isExcluded returns the exclude pattern that the given file matches, or an empty string if the
// file does not match any exclude patterns.
func (m *Globber) isExcluded(file string) (pattern string) {
	for _, pattern := range m.Exclude {
		if matched, _ := doublestar.Match(pattern, file); matched {
			return pattern
		}
	}

	return
}

// contains returns true if the given file is already present in the list of matched files, false otherwise.
func (m *Globber) contains(file string) bool {
	return slices.Contains(m.files, file)
}

// isExplicitlyIncluded returns true if the given file is considered to be explicitly included, which
// means the full pattern and the filename do not contain any glob characters.
func (m *Globber) isExplicitlyIncluded(file string) bool {
	noGlobs := !strings.Contains(file, "*")
	noGlobsInFilename := !strings.Contains(filepath.Base(file), "*")

	return noGlobs && noGlobsInFilename
}

// Match finds all files matching the given pattern and applies the exclusion options. After
// running this function, the matched files can be retrieved using the ListFiles method. Returns
// an error if the pattern fails to match.
func (m *Globber) Match(pattern string) (err error) {
	// Get all files that match the pattern
	var matches []string

	if matches, err = doublestar.FilepathGlob(pattern, doublestar.WithFilesOnly()); err != nil {
		return fmt.Errorf("failed to match pattern %q: %w", pattern, err)
	}

	for _, match := range matches {
		// Convert to absolute path
		match, _ = filepath.Abs(match)
		match = filepath.ToSlash(match)

		switch {
		// 1) Skip files that are already found
		case m.contains(match):
			m.Logger.Printf("<skipped> %q <already in matches>", match)
		// 2) If the file is explicitly included (i.e no glob pattern is used), then it should be included immediately.
		case m.isExplicitlyIncluded(pattern):
			m.Logger.Printf("<exception> %q <explicitly included>", match)
			m.files = append(m.files, match)
		case m.isExcluded(match) != "":
			m.Logger.Printf("<skipped> %q <matches exclude pattern> %q", match, m.isExcluded(match))
		case m.isBinary(match):
			m.Logger.Printf("<skipped> %q <detected as binary>", match)
		default:
			// Append the match to the matches slice
			m.files = append(m.files, match)
		}
	}

	return nil
}
