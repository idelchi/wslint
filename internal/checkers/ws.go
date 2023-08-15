package checkers

import (
	"errors"

	"github.com/idelchi/wslint/trailing"
)

// ErrHasTrailing is returned when there is trailing whitespace.
var ErrHasTrailing = errors.New("has trailing whitespace")

// Whitespace keeps track of trailing whitespaces.
type Whitespace struct {
	// Row(s) where the whitespace is found.
	rows []int
	// Error associated with the whitespace.
	error error
}

// Analyze determines if the line has trailing whitespace(s) and appends the row to the list of rows.
func (w *Whitespace) Analyze(line string, row int) {
	if trailing.Has(line) {
		w.rows = append(w.rows, row)
	}
}

// Finalize evaluates the correctness of trailing whitespace(s).
func (w *Whitespace) Finalize() {
	if w.rows != nil {
		w.error = ErrHasTrailing
	}
}

// Results returns the rows and error associated with the whitespace.
func (w *Whitespace) Results() ([]int, error) {
	return w.rows, w.error
}

// Stop returns 0.
func (w *Whitespace) Stop() int {
	return 0
}

// Fix removes trailing whitespace(s) from the line.
func (w *Whitespace) Fix(line string) string {
	return trailing.Trim(line)
}

// Info returns extra information.
func (w *Whitespace) Info() []string {
	return nil
}
