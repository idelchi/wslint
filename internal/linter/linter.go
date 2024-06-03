// Package linter provides a high-level interface to organize line-base formatting of a sequence of strings.
package linter

import (
	"log"

	"github.com/fatih/color"

	"github.com/idelchi/wslint/internal/checkers"
)

// Checker represents a line analyser.
type Checker interface {
	Format(lines []string) ([]string, []error)
}

// Linter represents a text linter.
type Linter struct {
	Name string
	// Checkers is a map with name and checker to use.
	Checkers map[string]Checker
	// Error contains the error, if any.
	Errors map[string][]error
	// The lines formatted.
	Lines []string
}

// InsertChecker adds a checker to the list of checkers in use.
func (l *Linter) InsertChecker(name string, c Checker) {
	l.Checkers[name] = c
}

// New creates a new linter, with the default checkers.
func New(name string) *Linter {
	defaultCheckers := map[string]Checker{
		"whitespace": checkers.Whitespace{},
		"blanks":     checkers.Blanks{},
	}

	return &Linter{
		Name:     name,
		Checkers: defaultCheckers,
		Errors:   make(map[string][]error),
	}
}

// HasCheckers returns true if the linter has checkers configured.
func (l *Linter) HasCheckers() bool {
	return l.Checkers != nil && len(l.Checkers) > 0
}

// HasIssues returns true if the linter has issues.
func (l *Linter) HasIssues() bool {
	return len(l.Errors) > 0
}

// Format returns the formatted string.
func (l *Linter) Format(lines []string) []string {
	if !l.HasCheckers() {
		panic("no checkers configured")
	}

	var errors []error
	for name, checker := range l.Checkers {
		if lines, errors = checker.Format(lines); len(errors) > 0 {
			l.Errors[name] = errors
		}
	}

	return lines
}

// Summary prints a summary of the file.
func (l *Linter) Summary() (ok bool) {
	// Use coloured output for emphasis
	filename := color.New(color.FgGreen, color.Bold).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()

	ok = !l.HasIssues()

	if !ok {
		log.Println(filename(l.Name))
		errors := l.Errors

		// Get all the keys in errors
		// checkers := maps.Keys(errors)

		for name, errors := range errors {
			log.Println("  - Errors detected: ", errorColor(name))

			for _, err := range errors {
				log.Printf("    - %s", errorColor(err))
			}
		}
	}

	return ok
}
