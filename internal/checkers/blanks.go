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

type Blanks struct{}

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

	slices.Reverse(rows)

	return
}

func (b Blanks) assert(rows []int) (errors []error) {
	switch blanks := len(rows); blanks {
	// one blank line at the end
	case 1:
		errors = nil
	// no blank lines at the end
	case 0:
		errors = append(errors, ErrTooFewBlanks)
	// more than one blank line at the end
	default:
		errors = append(errors, fmt.Errorf("%w: rows %v", ErrTooManyBlanks, rows))
	}

	return
}

func (b Blanks) format(lines []string, rows []int) []string {
	switch blanks := len(rows); blanks {
	// one blank line at the end
	case 1:
		return lines
	// no blank lines at the end
	case 0:
		return append(lines, "")
	// more than one blank line at the end
	default:
		return lines[:rows[1]]
	}
}

func (b Blanks) Format(lines []string) ([]string, []error) {
	rows := b.check(lines)
	errs := b.assert(rows)
	if len(errs) == 0 {
		return lines, errs
	}

	return b.format(lines, rows), errs
}
