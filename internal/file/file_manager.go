// Package file offers two main structures:
//   - Manager, which wraps a file and provides a simple interface for opening, closing and reading lines.
//   - Replacer, which wraps two Manager instances and provides a simple interface for safely replacing a file with
//     another.
//
// The intention is that the replacer is used by in-place formatters.
package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// closeFile closes a ReadWriteCloser and returns a error if it fails.
func closeFile(c io.ReadWriteCloser) (err error) {
	err = c.Close()

	switch {
	// If the file is already closed, return nil
	case errors.Is(err, os.ErrClosed):
		return nil
	case err != nil:
		return fmt.Errorf("failed to close: %w", err)
	default:
		return nil
	}
}

// Manager represents a wrapped management of file handling.
// Given a name, it can open, close and read lines from the file, until EOF.
type Manager struct {
	// Name is the name of the file to open. Should be a full or relative path.
	Name string
	// file is the file handle.
	file *os.File
	// reader is the buffered reader.
	reader *bufio.Reader
	// done is true if the file has been read to EOF.
	done bool
}

// Open the file for reading.
// Returns an error if the file doesn't exist.
func (f *Manager) Open() (err error) {
	f.file, err = os.Open(f.Name)

	if err != nil || f.file == nil {
		return fmt.Errorf("failed to open file %q: %w", f.Name, err)
	}

	f.reader = bufio.NewReader(f.file)

	return nil
}

// HasLines returns true if there are lines available to read.
func (f *Manager) HasLines() bool {
	return !f.done
}

// Next reads the next line from the file.
func (f *Manager) Next() (line string, err error) {
	line, err = f.reader.ReadString('\n')
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

// Close closes the file (if open).
// An optional error pointer can be passed in order to properly defer the error handling from the calling context.
func (f *Manager) Close(errPtr ...*error) (err error) {
	// Record the error
	err = closeFile(f.file)

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
