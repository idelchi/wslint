// Package trailing provides functions to check and remove trailing whitespace(s) from a string.
// Defines a whitespace according to the unicode.IsSpace function.
package trailing

import (
	"strings"
	"unicode"
)

// Has checks if a string has trailing whitespace, as defined by unicode.IsSpace.
// It returns true if the input string has trailing whitespace, false otherwise.
func Has(line string) bool {
	return line != strings.TrimRightFunc(line, unicode.IsSpace)
}

// Trim removes trailing whitespace, as defined by unicode.IsSpace, from a string.
// It returns a new string with the trailing whitespace removed.
func Trim(line string) string {
	return strings.TrimRightFunc(line, unicode.IsSpace)
}
