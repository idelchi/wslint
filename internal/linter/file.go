package linter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// File represents a wrapped management of file handling, using os.File and bufio.Reader.
// Given a name, it can open, close and read lines from the file, until EOF.
type File struct {
	// Name is the name of the file to open. Should be a full or relative path.
	Name string
	// File is the file handle.
	File *os.File
	// File is the buffered Readef.
	Reader *bufio.Reader
	// done is true if the file has been read to EOF.
	done bool
}

// NewFile creates a new File and opens it for reading.
func NewFile(name string) (File, error) {
	f := File{
		Name: name,
	}

	return f, f.Open()
}

func (f *File) Reset() (err error) {
	_, err = f.File.Seek(0, io.SeekStart)

	return
}

// Open the file for reading.
// Returns an error if the file doesn't exist.
func (f *File) Open() (err error) {
	// Open the file in read/write mode
	f.File, err = os.OpenFile(f.Name, os.O_RDWR, os.ModeAppend)

	if err != nil || f.File == nil {
		return fmt.Errorf("file manager failed to: %w", err)
	}

	f.Reader = bufio.NewReader(f.File)

	return
}

// HasLines returns true if there are lines available to read.
func (f *File) HasLines() bool {
	return !f.done
}

// Next reads the next line from the file.
func (f *File) Next() (line string, err error) {
	line, err = f.Reader.ReadString('\n')
	line = strings.TrimRight(strings.TrimRight(line, "\r\n"), "\n")

	switch {
	case errors.Is(err, io.EOF):
		f.done = true
		err = nil
	case err != nil:
		err = fmt.Errorf("failed to read line: %w", err)
	}

	return
}

// Close closes the file.
// Returns nil if the file is already closed.
func (f *File) Close() error {
	switch err := f.File.Close(); {
	// If the file is already closed, suppress the error and return nil
	case errors.Is(err, os.ErrClosed):
		return nil
	case err != nil:
		return fmt.Errorf("failed to close: %w", err)
	default:
		return nil
	}
}

// Write writes a line to the file.
func (f *File) Write(lines ...string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(f.File, line); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}

func (f *File) Rename(name string) error {
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	if err := os.Rename(f.Name, name); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	f.Name = name

	return f.Open()
}

type Renameable interface {
	Name() string
	Rename(string) error
	io.Closer
}

// Replace replaces two NamedCloser files.
func Replace(o, r Renameable) (err error) {
	name := o.Name()

	// Close the original file
	if err = o.Close(); err != nil {
		return fmt.Errorf("failed to replace %s: %w", name, err)
	}

	// Rename the replacement file
	if err = r.Rename(name); err != nil {
		return fmt.Errorf("failed to replace %s: %w", name, err)
	}

	return
}
