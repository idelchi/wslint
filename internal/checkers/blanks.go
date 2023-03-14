package checkers

import (
	"errors"
	"strings"
)

var (
	// ErrTooFewBlanks is returned when there are no blank lines at the end of the file.
	ErrTooFewBlanks = errors.New("no blank lines at the end of the file")
	// ErrTooManyBlanks is returned when there are more than one blank lines at the end of the file.
	ErrTooManyBlanks = errors.New("more than one blank line at the end of the file")
)

// Blanks keeps track of the number of blank lines at the end of the file.
type Blanks struct {
	// Rows are the rows (at the end) at which blank lines occur.
	Rows []int
	// Error associated with the number of blank lines.
	Error error
}

// Analyze determines whether the line is blank or not, and records its row number accordingly.
// It considers a line to be blank if it is empty after trimming the leading and trailing spaces.
func (b *Blanks) Analyze(line string, row int) {
	line = strings.TrimSpace(line)
	if line != "" {
		// Not a blank line, reset the record of blank lines.
		b.Rows = nil
	} else {
		// Blank line, record the row number.
		b.Rows = append(b.Rows, row)
	}
}

// Finalize evaluates the correctness of blank lines at the end of the file.
// If b.Rows is empty, there are no blank lines at the end of the file.
// If b.Rows has only one element, then there is one blank line at the end of the file.
// Any other value, means there are too many blank lines at the end of the file.
func (b *Blanks) Finalize() {
	switch blanks := len(b.Rows); blanks {
	// one blank line at the end
	case 1:
		b.Error = nil
	// no blank lines at the end
	case 0:
		b.Error = ErrTooFewBlanks
	// more than one blank line at the end
	default:
		b.Error = ErrTooManyBlanks
	}
}

// Results returns the rows at which blank lines occur, and the error associated with the number of blank lines.
func (b *Blanks) Results() ([]int, error) {
	return b.Rows, b.Error
}

// Stop returns the row at which the last useful blank line occurs.
func (b *Blanks) Stop() int {
	if b.Rows != nil && b.Error != nil {
		return b.Rows[0]
	}

	return 0
}

// Fix returns the line as is.
func (b *Blanks) Fix(line string) string {
	return line
}
