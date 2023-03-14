package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Manager represents a wrapped management of file handling.
type Manager struct {
	Name   string
	file   *os.File
	reader *bufio.Reader
	done   bool
}

// closeFile closes the ReadWriteCloser and returns the error.
func closeFile(c io.ReadWriteCloser) (err error) {
	err = c.Close()

	switch {
	case errors.Is(err, os.ErrClosed):
		err = nil
	case err != nil:
		err = fmt.Errorf("failed to close: %w", err)
	default:
		return nil
	}

	return
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
func (f *Manager) Close(errPtr ...*error) (err error) {
	err = closeFile(f.file)

	// If pointer syntax is used, set the error
	if len(errPtr) > 0 {
		errOuter := errPtr[0]
		// If the error is nil, set it
		// Else, wrap it.
		if *errOuter == nil {
			*errOuter = err
		} else {
			*errOuter = fmt.Errorf("%w, failed to close file: %w", *errOuter, err)
		}
	}

	return
}
