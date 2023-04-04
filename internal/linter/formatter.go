package linter

import (
	"fmt"
	"os"
	"path/filepath"
)

type Main = *Reader

// Formatter enables formatting of a file, by first writing to a temporary file.
// When the formatting is done, the original file is replaced with the temporary file.
type Writer struct {
	Main
	Shadow *Reader
}

// Write writes a line to the Shadow file.
func (f *Writer) Write(line ...string) error {
	return f.Shadow.Write(line...)
}

// Save applies the changes to the original file.
func (f *Writer) Save() error {
	return f.ReplaceWith(f.Shadow)
}

// Close closes both files.
func (f *Writer) Close() (err error) {
	if errClose := f.Main.Close(); errClose != nil {
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
func (f *Writer) Open() (err error) {
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

// Load loads both files.
func (f *Writer) Load(name string) (err error) {
	if errLoad := f.Main.Load(name); errLoad != nil {
		err = fmt.Errorf("%w", errLoad)
	}

	f.Shadow, err = CreateShadow(name)

	return
}

// NewWriter creates a new formatter for the given file.
func NewWriter(file *Reader) (formatter *Writer, err error) {
	formatter = &Writer{
		Main: file,
	}

	// formatter.Shadow, err = CreateShadow(file.Name)

	return
}

// CreateShadow creates a new shadow file for the given file.
func CreateShadow(name string) (shadow *Reader, err error) {
	// Create a replacement file to write the fixed file to
	tmpFile, err := os.CreateTemp(filepath.Dir(name), filepath.Base(name)+"-replacement-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	if shadow, err = NewReader(tmpFile.Name()); err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	return shadow, nil
}
