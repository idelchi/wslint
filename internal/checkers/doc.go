// Package checkers contains line analysis tooling.
// Each checker operates on a single line of text, while keeping track of which rows the errors occur.
// Some checkers require a full pass of the file, while others operate line by line.
// Example checkers:
// - Check for trailing whitespace
// - Check for trailing empty line at the end of a sequence of lines
package checkers
