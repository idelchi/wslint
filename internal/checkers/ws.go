package checkers

import (
	"errors"
	"fmt"

	"github.com/idelchi/wslint/trailing"
)

// TODO(Idelchi): Returns trailing whitespace on one-line documents with one blank line
// Need to restore the tests in this package.

// ErrHasTrailing is returned when there is trailing whitespace.
var ErrHasTrailing = errors.New("has trailing whitespace")

// Whitespace keeps track of trailing whitespaces.
type Whitespace struct{}

// Check identifies the lines that have trailing whitespaces.
func (w Whitespace) Check(lines []string) (rows []int) {
	for i, line := range lines {
		if trailing.Has(line) {
			rows = append(rows, i)
		}
	}

	return
}

// assert evaluates the correctness of trailing whitespaces.
// If rows is not empty, it means some lines have trailing whitespaces.
func (w Whitespace) assert(rows []int) (errors []error) {
	if len(rows) > 0 {
		// TODO(Idelchi): Would be clearer to the user if the row values are incremented by 1.
		errors = append(errors, fmt.Errorf("%w: on rows %v", ErrHasTrailing, rows))
	}

	return
}

// format removes trailing whitespaces from lines identified in rows.
func (w Whitespace) format(lines []string, rows []int) []string {
	for _, i := range rows {
		lines[i] = trailing.Trim(lines[i])
	}

	return lines
}

// Format checks the lines for trailing whitespaces, asserts any errors,
// and then formats the lines to remove those whitespaces.
func (w Whitespace) Format(lines []string) ([]string, []error) {
	rows := w.Check(lines)
	errs := w.assert(rows)

	if len(errs) == 0 {
		return lines, errs
	}

	return w.format(lines, rows), errs
}
