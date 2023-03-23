// Package linter provides functions to check and fix files for trailing whitespaces and blank lines.
// It offers three main components:
//   - Manager, which wraps a file and provides a simple interface for opening, closing and reading lines.
//   - Replacer, which wraps two Manager instances and provides a simple interface for safely replacing a file with
//     another. The intention is that the replacer is used by in-place formatters.
//   - Linter, which wraps a Manager and a list of checkers and provides a simple interface for linting and fixing
package linter
