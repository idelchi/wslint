// Package trailing provides functions to check and remove trailing whitespace(s) from a string.
// Defines a whitespace according to the unicode.IsSpace function.
package trailing

import (
	"strings"
	"unicode"
)

// Has checks if a string has trailing whitespace, as defined by unicode.IsSpace.
func Has(line string) bool {
	switch line {
	case "", strings.TrimRightFunc(line, unicode.IsSpace):
		return false
	default:
		return true
	}
}

// Trim removes trailing whitespace, as defined by unicode.IsSpace, from a string.
func Trim(line string) string {
	if !Has(line) {
		return line
	}

	return strings.TrimRightFunc(line, unicode.IsSpace)
}
