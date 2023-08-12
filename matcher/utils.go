package matcher

import (
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/godoc/util"
	"golang.org/x/tools/godoc/vfs"

	"github.com/bmatcuk/doublestar/v4"
)

// isBinary returns true if the given file is detected as a binary file, false otherwise.
func isBinary(file string) bool {
	fs := vfs.OS(filepath.Dir(file))

	return !util.IsTextFile(fs, filepath.Base(file))
}

// isExplicitlyIncluded returns true if the given file is considered to be explicitly included, which
// means the full pattern and the filename do not contain any glob characters.
func isExplicitlyIncluded(file string) bool {
	globsInPath := strings.Contains(file, "*")
	globsInName := strings.Contains(filepath.Base(file), "*")
	globsInExtension := strings.Contains(filepath.Ext(file), "*") || filepath.Ext(file) == ""

	return !globsInPath || (globsInName && !globsInExtension)
}

// isExcluded returns the exclude pattern that the given file matches, or an empty string if the
// file does not match any exclude patterns.
func isExcluded(file string, excludes []string) (pattern string) {
	for _, pattern := range excludes {
		if matched, _ := doublestar.Match(pattern, file); matched {
			return pattern
		}
	}

	return
}

// contains returns true if the given file is already present in the list of matched files, false otherwise.
func contains(file string, files []string) bool {
	return slices.Contains(files, file)
}
