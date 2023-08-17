package checkers

import (
	"errors"
	"fmt"
	"slices"

	"github.com/idelchi/wslint/stuttering"
)

// ErrStutter is returned when there is stuttering.
var ErrStutter = errors.New("stutters")

// Stutter keeps track of stuttering words.
type Stutter struct {
	Exceptions []string
}

func (s Stutter) check(lines []string) (rows []int, stutters map[int][]string) {
	stutters = make(map[int][]string)

	for row, line := range lines {
		if stuttering.Has(line) {
			words := stuttering.Find(line)

			for _, word := range words {
				if !slices.Contains(s.Exceptions, word) {
					stutters[row] = words
					rows = append(rows, row)
				}
			}
		}
	}
	return
}

func (s Stutter) assert(rows []int, stutters map[int][]string) (errors []error) {
	if len(rows) > 0 {
		for _, row := range rows {
			// TODO(Idelchi): Would be clearer to the user if the row values are incremented by 1.
			errors = append(errors, fmt.Errorf("%w: on line %d: words %v", ErrStutter, row, stutters[row]))
		}
	}
	return
}

func (s Stutter) format(lines []string, rows []int) []string {
	for _, i := range rows {
		lines[i] = stuttering.Trim(lines[i])
	}

	return lines
}

func (s Stutter) Format(lines []string) ([]string, []error) {
	rows, stutters := s.check(lines)
	errs := s.assert(rows, stutters)
	if len(errs) == 0 {
		return lines, errs
	}

	return s.format(lines, rows), errs
}
