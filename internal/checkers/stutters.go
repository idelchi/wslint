package checkers

import (
	"errors"
	"fmt"

	"github.com/idelchi/wslint/stuttering"
)

// ErrHasTrailing is returned when there is trailing stutter.
var ErrStutter = errors.New("stutters")

// Stutters keeps track of stuttering words.
type Stutters struct{}

func (s Stutters) check(lines []string) (rows []int, stutters map[int][]string) {
	stutters = make(map[int][]string)
	for i, line := range lines {
		if stuttering.Has(line) {
			words := stuttering.Find(line)
			stutters[i] = words
			rows = append(rows, i)
		}
	}
	return
}

func (s Stutters) assert(rows []int, stutters map[int][]string) (errors []error) {
	if len(rows) > 0 {
		for i, row := range rows {
			errors = append(errors, fmt.Errorf("%w: on line %d: words %v", ErrStutter, row, stutters[i]))
		}
	}
	return
}

func (s Stutters) format(lines []string, rows []int) []string {
	for _, i := range rows {
		lines[i] = stuttering.Trim(lines[i])
	}
	return lines
}

func (s Stutters) Format(lines []string) ([]string, []error) {
	rows, stutters := s.check(lines)
	errs := s.assert(rows, stutters)
	if len(errs) == 0 {
		return lines, errs
	}

	return s.format(lines, rows), errs
}
