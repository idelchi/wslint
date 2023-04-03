package linter

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/exp/slices"

	"github.com/idelchi/wslint/internal/checkers"
)

// Checker represents a line analyzer.
type Checker interface {
	Analyze(string, int)
	Finalize()
	Results() ([]int, error)
	Fix(string) string
	Stop() int
}

// Linter represents a file linter.
type Linter struct {
	// File to process.
	File Formatter
	// Checkers contains the checkers to be used.
	Checkers []Checker
	// Error contains the error, if any.
	Error error
	// Touched is a flag whether the file has been touched.
	Touched bool
}

// InsertChecker adds a checker to the list of checkers in use.
func (l *Linter) InsertChecker(c Checker) {
	l.Checkers = append(l.Checkers, c)
}

// NewLinter creates a new linter, with the default checkers.
func NewLinter(file *Formatter) *Linter {
	defaultCheckers := []Checker{
		&checkers.Whitespace{},
		&checkers.Blanks{},
	}

	return &Linter{
		File:     *file,
		Checkers: defaultCheckers,
	}
}

// HasCheckers returns true if the linter has checkers configured.
func (l *Linter) HasCheckers() bool {
	return l.Checkers != nil && len(l.Checkers) > 0
}

// Lint checks the file .
func (l *Linter) Lint() (err error) {
	if !l.HasCheckers() {
		return ErrNoCheckers
	}

	file := l.File

	file.Shadow = file.File

	if err := file.Open(); err != nil {
		return err
	}

	defer file.Close()

	for _, c := range l.Checkers {
		defer c.Finalize()
	}

	for row := 1; file.HasLines(); row++ {
		line, err := file.Next()
		if err != nil {
			return fmt.Errorf("error reading file: %w", err) //cover:ignore
		}

		for _, c := range l.Checkers {
			c.Analyze(line, row)
		}
	}

	return err
}

// ErrNoCheckers is returned when no checkers have been configured.
var ErrNoCheckers = errors.New("no checkers configured")

// checkErrs takes a list of errors and returns false if all errors are nil.
func checkErrs(errs ...error) bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}

	return false
}

// Fix fixes the file by removing trailing whitespaces and blank lines.
//
//nolint:cyclop,funlen
func (l *Linter) Fix() (err error) {
	if err = l.Lint(); err != nil {
		return err
	}

	// Array of errors. If any of the checkers have found errors, they will be added to this array.
	errs := make([]error, len(l.Checkers))

	// If none of the checkers have found any errors, there's nothing to fix and we can return.
	for i, c := range l.Checkers {
		_, errs[i] = c.Results()
	}

	// If there are no errors, return.
	if !checkErrs(errs...) {
		return nil
	}

	file := l.File

	file.PrepareForFixing()

	if err := file.Open(); err != nil {
		return err
	}

	defer file.Close()

	// For each checker, check if a stop row has been set.
	// This is used to stop the loop when the last error has been fixed.
	stops := make([]int, len(l.Checkers))

	for i, c := range l.Checkers {
		stops[i] = c.Stop()
	}

	// Sort the stop rows to get the highest one
	slices.Sort(stops)
	// Fetch the last stop row
	stop := stops[len(stops)-1]

	// // Write the fixed file to the temporary file
	for row := 1; file.HasLines(); row++ {
		if row == stop {
			break
		}

		line, err := file.Next()
		eof := !file.HasLines()

		if err != nil {
			return fmt.Errorf("error getting next line: %w", err) //cover:ignore
		}

		// If a line contains trailing whitespace, remove it
		// Not super efficient, but simpler. It becomes only a problem if the number of rows are huge.
		for _, c := range l.Checkers {
			rows, err := c.Results()
			if err != nil && slices.Contains(rows, row) {
				line = c.Fix(line)
			}
		}

		// This is an annoying edge case. If the last line is empty, the there will already
		// have been a newline written to the temporary file. If we don't check for this,
		// we will end up with two newlines at the end of the file.
		if eof && line == "" {
			break
		}

		if err := file.Write(line); err != nil {
			return fmt.Errorf("failed to copy line %d to temporary file: %w", row, err) //cover:ignore
		}
	}

	if err = file.Save(); err != nil {
		return //cover:ignore
	}

	l.Touched = true

	return err
}

// Summary prints a summary of the file.
func (l *Linter) Summary() (ok bool) {
	// If the file itself had an error, print it and return.
	if err := l.Error; err != nil {
		log.Println(l.File.Name)
		log.Println(err)

		return
	}

	ok = true

	messages := []string{}

	for _, c := range l.Checkers {
		rows, err := c.Results()
		if err != nil {
			messages = append(messages, fmt.Sprintf("- %v: at row(s): %v", err, rows))
			ok = false
		}
	}

	if !ok {
		log.Println(l.File.Name)

		for _, m := range messages {
			log.Println(m)
		}
	}

	if l.Touched {
		log.Printf("*** fixed ***")
	}

	return ok
}
