package checkers

import (
	"errors"

	"github.com/idelchi/wslint/internal/trailing"
)

// ErrHasTrailing is returned when there is trailing whitespace.
var ErrHasTrailing = errors.New("has trailing whitespace")

// Whitespace keeps track of trailing whitespaces.
type Whitespace struct {
	// Row(s) where the whitespace is found.
	Rows []int
	// Error associated with the whitespace.
	Error error
}

// Analyze determines if the line has trailing whitespace(s) and appends the row to the list of rows.
func (w *Whitespace) Analyze(line string, row int) {
	if trailing.Has(line) {
		w.Rows = append(w.Rows, row)
	}
}

// Finalize evaluates the correctness of trailing whitespace(s).
func (w *Whitespace) Finalize() {
	if w.Rows != nil {
		w.Error = ErrHasTrailing
	}
}

// Results returns the rows and error associated with the whitespace.
func (w *Whitespace) Results() ([]int, error) {
	return w.Rows, w.Error
}

// Stop returns 0.
func (w *Whitespace) Stop() int {
	return 0
}

// Fix removes trailing whitespace(s) from the line.
func (w *Whitespace) Fix(line string) string {
	return trailing.Trim(line)
}
