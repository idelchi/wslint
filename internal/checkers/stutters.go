package checkers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/idelchi/wslint/stuttering"
)

// ErrHasTrailing is returned when there is trailing stutter.
var ErrStutter = errors.New("stutters")

// Stutters keeps track of stuttering words.
type Stutters struct {
	// Row(s) where the stutter is found.
	rows []int
	// Error associated with the stutter.
	error error
	// Lines with stutter(s).
	lines []string
}

// Analyze determines if the line has trailing stutter(s) and appends the row to the list of rows.
func (s *Stutters) Analyze(line string, row int) {
	if stuttering.Has(line) {
		s.rows = append(s.rows, row)
		stutters := stuttering.Find(line)
		line = fmt.Sprintf("'%s' contains the stutters '%s'", strings.TrimSpace(line), strings.Join(stutters, ", "))
		s.lines = append(s.lines, line)
	}
}

// Finalize evaluates the correctness of trailing stutter(s).
func (s *Stutters) Finalize() {
	if s.rows != nil {
		s.error = ErrStutter
	}
}

// Results returns the rows and error associated with the stutter.
func (s *Stutters) Results() ([]int, error) {
	return s.rows, s.error
}

// Stop returns 0.
func (s *Stutters) Stop() int {
	return 0
}

// Fix removes stutter(s) from the line.
func (s *Stutters) Fix(line string) string {
	return stuttering.Trim(line)
}

// Info returns extra information about the stutter(s).
func (s *Stutters) Info() []string {
	return s.lines
}
