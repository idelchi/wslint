/*
Wslint formats files, removing trailing whitespaces and enforcing exactly one blank line at the end of the file.

Usage:

	wslint [flags] [path ...]

You can specify one or more glob patterns (or simply paths) as path arguments.

Use quotes to avoid shell expansion.

The flags are:

	-w		Format the files. Without this flag, wslint only performs linting.
	-a		Include hidden files and folders.
	-e		Exclude patterns, separated by commas (e.g., *.log,*.tmp).
	-j		Set the number of parallel jobs.
	-h		Print help information.
	-v		Print the version number.
	-d		Show debug output.
*/
package main
