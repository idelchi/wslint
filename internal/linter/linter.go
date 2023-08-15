package linter

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/fatih/color"

	"github.com/idelchi/wslint/internal/checkers"
)

// Checker represents a line analyser.
type Checker interface {
	Analyze(string, int)
	Finalize()
	Results() ([]int, error)
	Fix(string) string
	Stop() int
	Info() []string
}

// Writeable describes a file that can be written to.
type Writeable interface {
	Readable
	Write(...string) error
	Save() error
}

// Readable describes a file that can be read from.
type Readable interface {
	Open() error
	Close() error
	HasLines() bool
	Next() (string, error)
	Load(...string) error
}

// Linter represents a file linter.
type Linter struct {
	// Name contains the name of file to lint.
	Name string
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

// New creates a new linter, with the default checkers.
func New(name string) *Linter {
	defaultCheckers := []Checker{
		&checkers.Whitespace{},
		&checkers.Blanks{},
	}

	return &Linter{
		Name:     name,
		Checkers: defaultCheckers,
	}
}

// HasCheckers returns true if the linter has checkers configured.
func (l *Linter) HasCheckers() bool {
	return l.Checkers != nil && len(l.Checkers) > 0
}

// Lint checks the file.
func (l *Linter) Lint(file Readable) (err error) {
	if !l.HasCheckers() {
		return ErrNoCheckers
	}

	if err := file.Load(l.Name); err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}

	// Defer the file closing.
	defer func() {
		if errClose := file.Close(); err == nil && errClose != nil {
			err = errClose
		}
	}()

	for _, c := range l.Checkers {
		defer c.Finalize()
	}

	for row := 1; file.HasLines(); row++ {
		line, err := file.Next()
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		for _, c := range l.Checkers {
			c.Analyze(line, row)
		}
	}

	return err
}

// ErrNoCheckers is returned when no checkers have been configured.
var ErrNoCheckers = errors.New("no checkers configured")

// GetIssues returns a list of rows with issues and a list of errors.
func (l *Linter) GetIssues() (rows [][]int, errs []error, infos [][]string) {
	for _, c := range l.Checkers {
		row, err := c.Results()
		info := c.Info()
		if err != nil {
			rows = append(rows, row)
			errs = append(errs, err)
			infos = append(infos, info)
		}
	}

	return
}

// HasIssues takes a list of errors and returns false if all errors are nil.
func (l *Linter) HasIssues() bool {
	_, errors, _ := l.GetIssues()

	return len(errors) > 0
}

// Fix fixes the file by removing issues as defined in the added checkers.
//
//nolint:cyclop,funlen
func (l *Linter) Fix(file Writeable) (err error) {
	// TODO(Idelchi): What will happen if Lint() has not been called?
	if !l.HasCheckers() {
		return ErrNoCheckers
	}

	if err := file.Load(l.Name); err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}

	defer func() {
		if errClose := file.Close(); err == nil && errClose != nil {
			err = errClose
		}
	}()

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

	// Write the fixed file to the temporary file
	for row := 1; file.HasLines(); row++ {
		if row == stop {
			break
		}

		line, err := file.Next()
		eof := !file.HasLines()

		if err != nil {
			return fmt.Errorf("error getting next line: %w", err)
		}

		// If a line contains an issue, remove it.
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
			return fmt.Errorf("failed to copy line %d to temporary file: %w", row, err)
		}
	}

	if err = file.Save(); err != nil {
		err = fmt.Errorf("failed to save temporary file: %w", err)
	}

	l.Touched = true

	return err
}

// Summary prints a summary of the file.
func (l *Linter) Summary() (ok bool) {
	// Use coloured output for emphasis
	filename := color.New(color.FgGreen, color.Bold).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()
	fixedMessage := color.New(color.FgYellow).SprintFunc()

	// If the file itself had an error, print it and return.
	if err := l.Error; err != nil {
		log.Println(filename(l.Name))
		log.Printf("  %s", errorColor(err))

		return ok
	}

	ok = !l.HasIssues()

	if !ok {
		log.Println(filename(l.Name))

		rows, err, info := l.GetIssues()

		for i := range err {
			log.Printf("  - %s: at row(s): %v", errorColor(err[i]), rows[i])

			if info := info[i]; info != nil {
				for _, info := range info {
					log.Printf("    	%s", errorColor(info))
				}
			}
		}
	}

	if l.Touched {
		log.Println(fixedMessage("*** fixed ***"))
	}

	return ok
}
