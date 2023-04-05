package linter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Main is an alias for Reader, to be used as a base for other types.
type Main = *Reader

// Writer enables formatting of a file, by first writing to a temporary file.
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
	log.Printf("Saving changes to %q (shadowed by %q)", f.Name, f.Shadow.Name)

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
func (f *Writer) Load(filename ...string) (err error) {
	name := f.Main.Name

	if len(filename) > 0 {
		name = filename[0]
	}

	var (
		errMain   error
		errShadow error
	)

	if errMain = f.Main.Load(name); errMain != nil {
		err = fmt.Errorf("%w", errMain)
	}

	if f.Shadow, errShadow = CreateShadow(name); errShadow != nil {
		if errMain != nil {
			err = fmt.Errorf("%w\n%w", errMain, errShadow)
		} else {
			err = fmt.Errorf("%w", errShadow)
		}
	}

	return
}

// NewWriter creates a new formatter for the given file.
func NewWriter(file *Reader) (formatter *Writer, err error) {
	formatter = &Writer{
		Main: file,
	}

	return
}

// CreateShadow creates a new shadow file for the given file.
func CreateShadow(name string) (shadow *Reader, err error) {
	// Create a replacement file to write the fixed file to
	tmpFile, err := os.CreateTemp(filepath.Dir(name), filepath.Base(name)+"-replacement-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	// Immediately close it, because we're opening it using another reference.
	if err = tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	if shadow, err = NewReader(tmpFile.Name()); err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	return shadow, nil
}
