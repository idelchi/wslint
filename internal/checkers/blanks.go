package checkers

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	// ErrTooFewBlanks is returned when there are no blank lines at the end of the file.
	ErrTooFewBlanks = errors.New("no blank lines at the end of the file")
	// ErrTooManyBlanks is returned when there are more than one blank lines at the end of the file.
	ErrTooManyBlanks = errors.New("more than one blank line at the end of the file")
)

// Blanks is a checker that checks for trailing empty lines at the end of a sequence of lines.
// It returns an error if there are no blank lines at the end of the file or if there are more than one.
// It returns the formatted lines with the correct number of blank lines at the end of the file.
type Blanks struct{}

// check checks for trailing empty lines at the end of a sequence of lines.
// It returns the rows that are blank (at the end).
func (b Blanks) check(lines []string) (rows []int) {
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			// Blank line, record the row number.
			rows = append(rows, i)
		} else {
			// Not a blank line, stop the analysis.
			break
		}
	}

	// Reverse the slice to get the last non-blank entry as the first index.
	slices.Reverse(rows)

	return
}

// assert returns an error based on the number of blank lines at the end of the sequence of lines.
func (b Blanks) assert(rows []int) []error {
	switch blanks := len(rows); blanks {
	// no blank lines at the end
	case 0:
		return []error{ErrTooFewBlanks}
	// one blank line at the end
	case 1:
		return nil
	// more than one blank line at the end
	default:
		return []error{fmt.Errorf("%w: rows %v", ErrTooManyBlanks, rows)}
	}
}

// format returns the formatted lines with the correct number of blank lines at the end of the sequence of lines.
func (b Blanks) format(lines []string, rows []int) []string {
	switch blanks := len(rows); blanks {
	// no blank lines at the end
	case 0:
		return append(lines, "")
	// one blank line at the end
	case 1:
		return lines
	// more than one blank line at the end
	default:
		return lines[:rows[1]]
	}
}

// Format checks the correctness of the sequence of lines in terms of blank lines at the end,
// applies the formatting if needed and returns the formatted lines along with the errors.
func (b Blanks) Format(lines []string) ([]string, []error) {
	rows := b.check(lines)
	errs := b.assert(rows)

	if len(errs) == 0 {
		return lines, errs
	}

	return b.format(lines, rows), errs
}
