/*
Wslint formats files, removing trailing whitespaces and enforcing exactly one blank line at the end of the file.

Usage:

	wslint [flags] [path ...]

Paths can be specified as one or more glob patterns or simple file paths.

Enclose the path arguments in quotes to prevent shell expansion.

The flags are:

	-w		Format the files. Without this flag, wslint only performs linting.
	-a		Include hidden files and folders.
	-e		Exclude patterns, separated by commas (e.g., *.log,*.tmp).
	-j		Set the number of parallel jobs.
	-h		Print help information.
	-v		Print the version number.
	-d		Show debug output.
	-q		Suppress messages.
*/
package main
