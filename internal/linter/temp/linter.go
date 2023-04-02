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

// type File interface {
// 	Name() string
// 	Read() ([]string, error)
// 	Open() error
// 	Close(*error) error
// 	Write(string) error
// 	Next() (string, error)
// 	HasLines() bool
// 	Finish() error
// }

// Linter represents a file linter.
type Linter struct {
	// Name of the file to lint.
	Name string
	// Checkers contains the checkers to be used.
	Checkers []Checker
	// Error contains the error, if any.
	Error error
	// Touched is a flag whether the file has been touched.
	Touched bool
	File    File
}

// InsertChecker adds a checker to the list of checkers in use.
func (l *Linter) InsertChecker(c Checker) {
	l.Checkers = append(l.Checkers, c)
}

// NewLinter creates a new linter, with the default checkers.
func NewLinter(file File) *Linter {
	defaultCheckers := []Checker{
		&checkers.Whitespace{},
		&checkers.Blanks{},
	}

	return &Linter{
		File:     file,
		Checkers: defaultCheckers,
	}
}

func (l *Linter) HasCheckers() bool {
	return l.Checkers != nil && len(l.Checkers) > 0
}

// Lint checks the file .
func (l *Linter) Lint() (err error) {
	if !l.HasCheckers() {
		return ErrNoCheckers
	}

	file := l.File

	defer file.Close(&err) //nolint:errcheck // The error is checked, with the &err parameter.

	// Open the file
	if err = file.Open(); err != nil {
		return fmt.Errorf("failed to lint: %w", err)
	}

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

// Fix fixes the file by removing trailing whitespaces and blank lines.
//
//nolint:cyclop,funlen
func (l *Linter) Fix() (err error) {
	if err = l.Lint(); err != nil {
		return err
	}

	// If none of the checkers have found any errors, there's nothing to fix and we can return.
	for _, c := range l.Checkers {
		_, errC := c.Results()
		if errC != nil {
			err = errC
		}
	}

	// If there are no errors, return.
	if err == nil {
		return nil
	}

	file := l.File

	defer file.Close(&err) //nolint:errcheck // The error is checked, with the &err parameter.
	// Open the file
	if err = file.Open(); err != nil {
		return err //cover:ignore
	}

	// For each checker, check if a stop row has been set.
	// This is used to stop the loop when the last error has been fixed.
	stops := make([]int, 0, len(l.Checkers))

	for _, c := range l.Checkers {
		stops = append(stops, c.Stop())
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

	if err = file.Finish(); err != nil {
		return //cover:ignore
	}

	l.Touched = true

	return err
}

// Summary prints a summary of the file.
func (l *Linter) Summary() (ok bool) {
	// If the file itself had an error, print it and return.
	if err := l.Error; err != nil {
		log.Println(l.Name)
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
		log.Println(l.Name)

		for _, m := range messages {
			log.Println(m)
		}
	}

	if l.Touched {
		log.Printf("*** fixed ***")
	}

	return ok
}
