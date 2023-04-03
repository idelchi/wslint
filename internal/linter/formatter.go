package linter

import (
	"fmt"
	"os"
	"path/filepath"
)

type Main = *File

// Formatter enables formatting of a file, by first writing to a temporary file.
// When the formatting is done, the original file is replaced with the temporary file.
type Formatter struct {
	Main
	Shadow *File
}

// Write writes a line to the Shadow file.
func (f *Formatter) Write(line ...string) error {
	return f.Shadow.Write(line...)
}

// Save applies the changes to the original file.
func (f *Formatter) Save() error {
	return f.ReplaceWith(f.Shadow)
}

// Close closes both files.
func (f *Formatter) Close() (err error) {
	if errClose := f.File.Close(); errClose != nil {
		err = fmt.Errorf("%w", errClose)
	}

	if errClose := f.Shadow.Close(); errClose != nil {
		if err != nil {
			err = fmt.Errorf("%w\n%w", err, errClose)
		} else {
			err = fmt.Errorf("%w", errClose)
		}
	}

	return
}

// Open opens both files.
func (f *Formatter) Open() (err error) {
	if errOpen := f.Main.Open(); errOpen != nil {
		err = fmt.Errorf("%w", errOpen)
	}

	if errOpen := f.Shadow.Open(); errOpen != nil {
		if err != nil {
			err = fmt.Errorf("%w\n%w", err, errOpen)
		} else {
			err = fmt.Errorf("%w", errOpen)
		}
	}

	return
}

// Cleanup removes the Shadow file.
func (f *Formatter) Cleanup() error {
	return f.Shadow.Delete()
}

// NewFormatter creates a new formatter for the given file.
func NewFormatter(file *File) (formatter *Formatter, err error) {
	formatter = &Formatter{
		Main: file,
	}

	formatter.Shadow, err = CreateShadow(file.Name)

	return
}

// CreateShadow creates a new shadow file for the given file.
func CreateShadow(name string) (shadow *File, err error) {
	// Create a replacement file to write the fixed file to
	tmpFile, err := os.CreateTemp(filepath.Dir(name), filepath.Base(name)+"-replacement-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	if shadow, err = NewFile(tmpFile.Name()); err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	return shadow, nil
}
