package linter

import (
	"fmt"
	"os"
	"path/filepath"
)

// Replacer represents a file replacer.
// It contains a Reader for the original file and the replacement file.
type Replacer struct {
	Original    Reader
	Replacement Reader
}

// Close closes both files (if open).
// It uses a pointer to the error to allow the caller to defer the call.
func (r *Replacer) Close(errPtr ...*error) (err error) {
	err1 := r.Original.Close()
	err2 := r.Replacement.Close()

	switch {
	case err1 != nil && err2 != nil:
		err = fmt.Errorf("\n%w\n%w", err1, err2)
	case err1 != nil:
		err = fmt.Errorf("original file: %w", err1)
	case err2 != nil:
		err = fmt.Errorf("replacement file: %w", err2)
	}

	// If pointer syntax is used, set the error
	if len(errPtr) > 0 {
		errOuter := errPtr[0]
		// If the error is nil, set it
		// Else, wrap it.
		if *errOuter == nil {
			*errOuter = err
		} else {
			*errOuter = fmt.Errorf("%w: %w", *errOuter, err)
		}
	}

	return
}

// NewReplacer creates a new file replacer.
func NewReplacer(name Reader) (Replacer, error) {
	r := Replacer{Original: name}

	return r, r.Setup()
}

// Setup sets up the replacement file.
func (r *Replacer) Setup() (err error) {
	name := r.Original.Name

	// Get the file parent directory
	parentDir := filepath.Dir(name)
	// Get the file name
	fileName := filepath.Base(name)

	// Create a replacement file to write the fixed file to
	tmpFile, err := os.CreateTemp(parentDir, fileName+"-replacement-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	r.Replacement = Reader{
		Name:   tmpFile.Name(),
		file:   tmpFile,
		reader: nil,
	}

	return
}

// Replace replaces the original file with the replacement file.
func (r *Replacer) Replace() (err error) {
	// Must close both files before replacing.
	if err = r.Close(); err != nil {
		return fmt.Errorf("failed to replace files: %w", err)
	}

	// Rename the temporary file to the original file
	if err := os.Rename(r.Replacement.Name, r.Original.Name); err != nil {
		return fmt.Errorf("failed to rename file %q to %q: %w", r.Replacement.Name, r.Original.Name, err)
	}

	return nil
}
