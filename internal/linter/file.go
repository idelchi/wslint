package linter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type ReadWriteSeekerCloser interface {
	io.ReadWriteSeeker
	io.Closer
}

type ReadStringer interface {
	ReadString(delim byte) (string, error)
}

// File represents a wrapped management of file handling, using os.File and bufio.Reader.
// Given a name, it can open, close and read lines from the file, until EOF.
type Reader struct {
	// Name is the name of the file to open. Should be a full or relative path.
	Name string
	// File is the file handle.
	file ReadWriteSeekerCloser
	// File is the buffered reader.
	buffer ReadStringer
	// done is true if the file has been read to EOF.
	done bool
}

// NewReader opens a file for reading (and writing).
func NewReader(name string) (*Reader, error) {
	f := &Reader{
		Name: name,
	}

	return f, f.Open()
}

// Open the file for reading.
// Returns an error if the file doesn't exist.
func (f *Reader) Open() (err error) {
	// Open the file in read/write mode
	f.file, err = os.OpenFile(f.Name, os.O_RDWR, os.ModeAppend)

	if err != nil || f.file == nil {
		return fmt.Errorf("file manager failed to: %w", err)
	}

	return f.Reset()
}

// Close closes the file.
// Returns nil if the file is already closed.
func (f *Reader) Close() error {
	switch err := f.file.Close(); {
	// If the file is already closed, suppress the error and return nil
	case errors.Is(err, os.ErrClosed):
		return nil
	case err != nil:
		return fmt.Errorf("failed to close: %w", err)
	default:
		return nil
	}
}

// Next reads the next line from the file.
func (f *Reader) Next() (line string, err error) {
	line, err = f.buffer.ReadString('\n')
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

// HasLines returns true if there are lines available to read.
func (f *Reader) HasLines() bool {
	return !f.done
}

// Reset resets the file to the beginning and assigns a fresh reader.
func (f *Reader) Reset() error {
	f.buffer = bufio.NewReader(f.file)
	f.done = false
	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file: %w", err)
	}

	return nil
}

// Write writes a line to the file.
func (f *Reader) Write(lines ...string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(f.file, line); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}

func (f *Reader) Rename(name string) error {
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	if err := os.Rename(f.Name, name); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	f.Name = name

	return nil
}

// ReplaceWith replaces the current file with the given file.
func (f *Reader) ReplaceWith(replacement *Reader) (err error) {
	name := f.Name

	// Close the original
	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to replace %s: %w", name, err)
	}

	return replacement.Rename(name)
}

// Save simply closes the file, since all writes are done in place.
func (f *Reader) Save() error {
	return f.Close()
}

// Load opens the file.
// Returns an error if the file doesn't exist.
func (f *Reader) Load(name string) error {
	f.Name = name

	return f.Open()
}
